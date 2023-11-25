package main

import (
	"database/sql"
	"errors"
	"net/http"
	"sort"
	"strconv"

	"github.com/labstack/echo/v4"
)

type LivestreamStatistics struct {
	Rank           int64 `json:"rank"`
	ViewersCount   int64 `json:"viewers_count"`
	TotalReactions int64 `json:"total_reactions"`
	TotalReports   int64 `json:"total_reports"`
	MaxTip         int64 `json:"max_tip"`
}

type LivestreamRankingEntry struct {
	LivestreamID int64
	Score        int64
}
type LivestreamRanking []LivestreamRankingEntry

func (r LivestreamRanking) Len() int      { return len(r) }
func (r LivestreamRanking) Swap(i, j int) { r[i], r[j] = r[j], r[i] }
func (r LivestreamRanking) Less(i, j int) bool {
	if r[i].Score == r[j].Score {
		return r[i].LivestreamID < r[j].LivestreamID
	} else {
		return r[i].Score < r[j].Score
	}
}

type UserStatistics struct {
	Rank              int64  `json:"rank"`
	ViewersCount      int64  `json:"viewers_count"`
	TotalReactions    int64  `json:"total_reactions"`
	TotalLivecomments int64  `json:"total_livecomments"`
	TotalTip          int64  `json:"total_tip"`
	FavoriteEmoji     string `json:"favorite_emoji"`
}

type UserRankingEntry struct {
	Username string `db:username`
	Score    int64  `db:score`
}
type UserRanking []UserRankingEntry

func (r UserRanking) Len() int      { return len(r) }
func (r UserRanking) Swap(i, j int) { r[i], r[j] = r[j], r[i] }
func (r UserRanking) Less(i, j int) bool {
	if r[i].Score == r[j].Score {
		return r[i].Username < r[j].Username
	} else {
		return r[i].Score < r[j].Score
	}
}

func getUserStatisticsHandler(c echo.Context) error {
	ctx := c.Request().Context()

	if err := verifyUserSession(c); err != nil {
		// echo.NewHTTPErrorが返っているのでそのまま出力
		return err
	}

	username := c.Param("username")
	// ユーザごとに、紐づく配信について、累計リアクション数、累計ライブコメント数、累計売上金額を算出
	// また、現在の合計視聴者数もだす

	tx, err := dbConn.BeginTxx(ctx, nil)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to begin transaction: "+err.Error())
	}
	defer tx.Rollback()

	var user UserModel
	if err := tx.GetContext(ctx, &user, "SELECT * FROM users WHERE name = ?", username); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return echo.NewHTTPError(http.StatusBadRequest, "not found user that has the given username")
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to get user: "+err.Error())
		}
	}

	// ランク算出
	// var users []*UserModel
	// if err := tx.SelectContext(ctx, &users, "SELECT * FROM users"); err != nil {
	// 	return echo.NewHTTPError(http.StatusInternalServerError, "failed to get users: "+err.Error())
	// }

	var ranking []*UserRankingEntry
	// select r.id as Username, (r.count + t.count) as Score from (SELECT u.id, COUNT(*) as count FROM users u INNER JOIN livestreams l ON l.user_id = u.id INNER JOIN reactions r ON r.livestream_id = l.id group by u.id) as r inner join (SELECT u.id, IFNULL(SUM(l2.tip), 0) as count FROM users u INNER JOIN livestreams l ON l.user_id = u.id INNER JOIN livecomments l2 ON l2.livestream_id = l.id group by u.id) as t on r.id = t.id order by Score, Username;
	query := `select r.Name as username, (r.count + t.count) as Score from (SELECT u.id, u.Name, COUNT(*) as count FROM users u INNER JOIN livestreams l ON l.user_id = u.id INNER JOIN reactions r ON r.livestream_id = l.id group by u.id) as r inner join (SELECT u.id, IFNULL(SUM(l2.tip), 0) as count FROM users u INNER JOIN livestreams l ON l.user_id = u.id INNER JOIN livecomments l2 ON l2.livestream_id = l.id group by u.id) as t on r.id = t.id order by score, username`
	if err := tx.GetContext(ctx, &ranking, query); err != nil && !errors.Is(err, sql.ErrNoRows) {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to count reactions: "+err.Error())
	}
	//(SELECT u.id as id, COUNT(*) as count FROM users u INNER JOIN livestreams l ON l.user_id = u.id INNER JOIN reactions r ON r.livestream_id = l.id group by u.id) as r;
	// (SELECT u.id as id, IFNULL(SUM(l2.tip), 0) as count FROM users u INNER JOIN livestreams l ON l.user_id = u.id	INNER JOIN livecomments l2 ON l2.livestream_id = l.id) as t;
	// select r.Name as username, (r.count + t.count) as Score from (SELECT u.id, u.Name, COUNT(*) as count FROM users u INNER JOIN livestreams l ON l.user_id = u.id INNER JOIN reactions r ON r.livestream_id = l.id group by u.id) as r inner join (SELECT u.id, IFNULL(SUM(l2.tip), 0) as count FROM users u INNER JOIN livestreams l ON l.user_id = u.id INNER JOIN livecomments l2 ON l2.livestream_id = l.id group by u.id) as t on r.id = t.id order by score, username;
	// for _, user := range users {
	// 	var reactions int64
	// 	// SELECT u.id, COUNT(*) as count FROM users u INNER JOIN livestreams l ON l.user_id = u.id INNER JOIN reactions r ON r.livestream_id = l.id group by u.id;
	// 	query := `
	// SELECT COUNT(*) FROM users u
	// INNER JOIN livestreams l ON l.user_id = u.id
	// INNER JOIN reactions r ON r.livestream_id = l.id
	// 	WHERE u.id = ?`
	// 	if err := tx.GetContext(ctx, &reactions, query, user.ID); err != nil && !errors.Is(err, sql.ErrNoRows) {
	// 		return echo.NewHTTPError(http.StatusInternalServerError, "failed to count reactions: "+err.Error())
	// 	}

	// 	var tips int64
	// 	query = `
	// 	SELECT IFNULL(SUM(l2.tip), 0) FROM users u
	// 	INNER JOIN livestreams l ON l.user_id = u.id
	// 	INNER JOIN livecomments l2 ON l2.livestream_id = l.id
	// 	WHERE u.id = ?`
	// 	if err := tx.GetContext(ctx, &tips, query, user.ID); err != nil && !errors.Is(err, sql.ErrNoRows) {
	// 		return echo.NewHTTPError(http.StatusInternalServerError, "failed to count tips: "+err.Error())
	// 	}

	// 	score := reactions + tips
	// 	ranking = append(ranking, UserRankingEntry{
	// 		Username: user.Name,
	// 		Score:    score,
	// 	})
	// }
	// sort.Sort(ranking)

	var rank int64 = 1
	for i := len(ranking) - 1; i >= 0; i-- {
		entry := ranking[i]
		if entry.Username == username {
			break
		}
		rank++
	}

	// リアクション数
	var totalReactions int64
	query = `SELECT COUNT(*) FROM users u 
    INNER JOIN livestreams l ON l.user_id = u.id 
    INNER JOIN reactions r ON r.livestream_id = l.id
    WHERE u.name = ?
	`
	if err := tx.GetContext(ctx, &totalReactions, query, username); err != nil && !errors.Is(err, sql.ErrNoRows) {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to count total reactions: "+err.Error())
	}

	// ライブコメント数、チップ合計
	var totalLivecomments int64
	var totalTip int64

	var total *LiveCommentTotal
	// mysql> select sum(subquery.count) as total_count, sum(subquery.tip_sum) as total_tip from (select livestreams.id as id, COUNT(*) as count, IFNULL(SUM(tip),0) as tip_sum from  livestreams as livestreams left join livecomments on livecomments.livestream_id = livestreams.id group by livestreams.id) as subquery;
	// select IFNULL(sum(subquery.count),0) as total_count, IFNULL(sum(subquery.tip_sum),0) as total_tip from (select livestreams.id as id, COUNT(*) as count, IFNULL(SUM(tip),0) as tip_sum from  livestreams as livestreams left join livecomments on livecomments.livestream_id = livestreams.id where livestreams.user_id = "test001" group by livestreams.id) as subquery;
	query = "select IFNULL(sum(subquery.count),0) as total_count, IFNULL(sum(subquery.tip_sum),0) as total_tip from (select livestreams.id as id, COUNT(*) as count, IFNULL(SUM(tip),0) as tip_sum from  livestreams as livestreams left join livecomments on livecomments.livestream_id = livestreams.id where livestreams.user_id = ? group by livestreams.id) as subquery"
	if err := tx.GetContext(ctx, &total, query, user.ID); err != nil && !errors.Is(err, sql.ErrNoRows) {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get livecomments: "+err.Error())
	}
	totalLivecomments = total.totalComment
	totalTip = total.totalTip
	// for _, livestream := range livestreams {
	// 	var livecomments []*LivecommentModel
	// 	if err := tx.SelectContext(ctx, &livecomments, "SELECT * FROM livecomments WHERE livestream_id = ?", livestream.ID); err != nil && !errors.Is(err, sql.ErrNoRows) {
	// 		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get livecomments: "+err.Error())
	// 	}

	// 	for _, livecomment := range livecomments {
	// 		totalTip += livecomment.Tip
	// 		totalLivecomments++
	// 	}
	// }

	query = "select IFNULL(SUM(view_count),0)  from (select * from livestreams where user_id = ?) as live inner join (select livestream_id, COUNT(*) as view_count from livestream_viewers_history group by livestream_id) as history on history.livestream_id = live.id"
	// 合計視聴者数
	var viewersCount int64
	if err := tx.GetContext(ctx, &viewersCount, "query", user.ID); err != nil && !errors.Is(err, sql.ErrNoRows) {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get livestream_view_history: "+err.Error())
	}

	// お気に入り絵文字
	var favoriteEmoji string
	query = `
	SELECT r.emoji_name
	FROM users u
	INNER JOIN livestreams l ON l.user_id = u.id
	INNER JOIN reactions r ON r.livestream_id = l.id
	WHERE u.name = ?
	GROUP BY emoji_name
	ORDER BY COUNT(*) DESC, emoji_name DESC
	LIMIT 1
	`
	if err := tx.GetContext(ctx, &favoriteEmoji, query, username); err != nil && !errors.Is(err, sql.ErrNoRows) {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to find favorite emoji: "+err.Error())
	}

	stats := UserStatistics{
		Rank:              rank,
		ViewersCount:      viewersCount,
		TotalReactions:    totalReactions,
		TotalLivecomments: totalLivecomments,
		TotalTip:          totalTip,
		FavoriteEmoji:     favoriteEmoji,
	}
	return c.JSON(http.StatusOK, stats)
}

func getLivestreamStatisticsHandler(c echo.Context) error {
	ctx := c.Request().Context()

	if err := verifyUserSession(c); err != nil {
		return err
	}

	id, err := strconv.Atoi(c.Param("livestream_id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "livestream_id in path must be integer")
	}
	livestreamID := int64(id)

	tx, err := dbConn.BeginTxx(ctx, nil)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to begin transaction: "+err.Error())
	}
	defer tx.Rollback()

	var livestream LivestreamModel
	if err := tx.GetContext(ctx, &livestream, "SELECT * FROM livestreams WHERE id = ?", livestreamID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return echo.NewHTTPError(http.StatusBadRequest, "cannot get stats of not found livestream")
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to get livestream: "+err.Error())
		}
	}

	var livestreams []*LivestreamModel
	if err := tx.SelectContext(ctx, &livestreams, "SELECT * FROM livestreams"); err != nil && !errors.Is(err, sql.ErrNoRows) {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get livestreams: "+err.Error())
	}

	// ランク算出
	var ranking LivestreamRanking
	for _, livestream := range livestreams {
		var reactions int64
		if err := tx.GetContext(ctx, &reactions, "SELECT COUNT(*) FROM livestreams l INNER JOIN reactions r ON l.id = r.livestream_id WHERE l.id = ?", livestream.ID); err != nil && !errors.Is(err, sql.ErrNoRows) {
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to count reactions: "+err.Error())
		}

		var totalTips int64
		if err := tx.GetContext(ctx, &totalTips, "SELECT IFNULL(SUM(l2.tip), 0) FROM livestreams l INNER JOIN livecomments l2 ON l.id = l2.livestream_id WHERE l.id = ?", livestream.ID); err != nil && !errors.Is(err, sql.ErrNoRows) {
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to count tips: "+err.Error())
		}

		score := reactions + totalTips
		ranking = append(ranking, LivestreamRankingEntry{
			LivestreamID: livestream.ID,
			Score:        score,
		})
	}
	sort.Sort(ranking)

	var rank int64 = 1
	for i := len(ranking) - 1; i >= 0; i-- {
		entry := ranking[i]
		if entry.LivestreamID == livestreamID {
			break
		}
		rank++
	}

	// 視聴者数算出
	var viewersCount int64
	if err := tx.GetContext(ctx, &viewersCount, `SELECT COUNT(*) FROM livestreams l INNER JOIN livestream_viewers_history h ON h.livestream_id = l.id WHERE l.id = ?`, livestreamID); err != nil && !errors.Is(err, sql.ErrNoRows) {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to count livestream viewers: "+err.Error())
	}

	// 最大チップ額
	var maxTip int64
	if err := tx.GetContext(ctx, &maxTip, `SELECT IFNULL(MAX(tip), 0) FROM livestreams l INNER JOIN livecomments l2 ON l2.livestream_id = l.id WHERE l.id = ?`, livestreamID); err != nil && !errors.Is(err, sql.ErrNoRows) {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to find maximum tip livecomment: "+err.Error())
	}

	// リアクション数
	var totalReactions int64
	if err := tx.GetContext(ctx, &totalReactions, "SELECT COUNT(*) FROM livestreams l INNER JOIN reactions r ON r.livestream_id = l.id WHERE l.id = ?", livestreamID); err != nil && !errors.Is(err, sql.ErrNoRows) {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to count total reactions: "+err.Error())
	}

	// スパム報告数
	var totalReports int64
	if err := tx.GetContext(ctx, &totalReports, `SELECT COUNT(*) FROM livestreams l INNER JOIN livecomment_reports r ON r.livestream_id = l.id WHERE l.id = ?`, livestreamID); err != nil && !errors.Is(err, sql.ErrNoRows) {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to count total spam reports: "+err.Error())
	}

	if err := tx.Commit(); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to commit: "+err.Error())
	}

	return c.JSON(http.StatusOK, LivestreamStatistics{
		Rank:           rank,
		ViewersCount:   viewersCount,
		MaxTip:         maxTip,
		TotalReactions: totalReactions,
		TotalReports:   totalReports,
	})
}
