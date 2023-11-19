scp でサーバー　 → ローカルにダウンロード
(isucon 配下) scp -ri "isucon12-allinone.pem" isucon@ec2-3-112-242-140.ap-northeast-1.compute.amazonaws.com:/home/isucon/webapp ./
(etc/とか) scp -ri "isucon12-allinone.pem" ubuntu@ec2-3-112-242-140.ap-northeast-1.compute.amazonaws.com:/etc/nginx/nginx.conf ./

ローカル　 → サーバー　に アップロード
sudo scp -i "isucon12-allinone.pem" -r ./webapp isucon@ec2-3-112-242-140.ap-northeast-1.compute.amazonaws.com:/home/isucon/

scp でセキュアな場所（sudo 必要）にアップロードしたい時(nginx.conf)
scp -i "isucon12-allinone.pem" -r ./nginx.conf ubuntu@ec2-3-112-242-140.ap-northeast-1.compute.amazonaws.com:/home/ubuntu/nginx.conf
ssh -i "isucon12-allinone.pem" ubuntu@ec2-3-112-242-140.ap-northeast-1.compute.amazonaws.com "sudo mv /home/ubuntu/nginx.conf /etc/nginx/nginx.conf"

のように２段階でやる必要あり

./scp_transfer.sh path/to/config.json
