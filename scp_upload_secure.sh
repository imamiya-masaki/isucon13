#!/bin/bash

# JSONファイルのパス
json_file="./shellConfig.json"

# 引数からターゲットのキーを取得
target=$1

# 引数がない場合はエラーメッセージを表示して終了
if [ -z "$json_file" ]; then
    echo "Usage: $0 <json_file>"
    exit 1
fi

# JSONファイルの内容を変数に読み込む
isucon_pem=$(jq -r ".isucon_pem" "$json_file")
isucon_user=$(jq -r ".isucon_user" "$json_file")
admin_user=$(jq -r ".admin_user" "$json_file")
host=$(jq -r ".host" "$json_file")
target_is_secure=$(jq ".file_mapper.$target.isSecure" "$json_file")
target_local_path=$(jq -r ".file_mapper.$target.path.local" "$json_file")
target_remote_path=$(jq -r ".file_mapper.$target.path.remote" "$json_file")


# 変数の中身を確認
echo "target: $target"
echo "isucon_pem: $isucon_pem"
echo "isucon_user: $isucon_user"
echo "admin_user: $admin_user"
echo "host: $host"
echo "target_is_secure: $target_is_secure"
echo "target_local_path: $target_local_path"
echo "target_remote_path: $target_remote_path"

 # target_is_secureがnullの場合はエラーメッセージを表示して終了
if [ "$target_is_secure" == null ]; then
    echo "Error: $target is null"
    exit 1
fi

# target_is_secureがtrueの場合はscpコマンドを実行, falseの場合はsshコマンドを実行
if [ "$target_is_secure" = true ]; then
    echo "sshコマンドを実行します"
    # targetの先頭の./を削除したtarget_dirを作成
    target_dir=$(echo "$target_local_path" | sed -e 's/^\.\///g')
    echo target_dir: $target_dir
    echo "scp -i $isucon_pem -r $target_local_path $isucon_user@$host:/home/$admin_user/$target_dir"
    echo "ssh -i $isucon_pem $admin_user@$host sudo mv /home/$admin_user/$target_dir $target_remote_path"
    scp -i "$isucon_pem" -r "$target_local_path" "$isucon_user@$host:/home/$admin_user/$target_dir"
    ssh -i "$isucon_pem" "$admin_user@$host" "sudo mv /home/$admin_user/$target_dir $target_remote_path"
else
    echo "scpコマンドを実行します"
    echo "scp -i $isucon_pem -r $target_local_path $isucon_user@$host:$target_remote_path"
    scp -i "$isucon_pem" -r "$target_local_path" "$isucon_user@$host:$target_remote_path"
fi