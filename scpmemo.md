# download server1
## mysql
```
scp -ri "~/.ssh/id_rsa" isucon@ec2-3-112-135-213.ap-northeast-1.compute.amazonaws.com:/etc/mysql/mysql.conf.d/mysqld.cnf ./server1/conf/
```
## nginx
```
scp -ri "~/.ssh/id_rsa" isucon@ec2-3-112-135-213.ap-northeast-1.compute.amazonaws.com:/etc/nginx/nginx.conf ./server1/conf/cd 
```
## slowquerydigest
```
```
scp -ri "~/.ssh/id_rsa" isucon@ec2-3-112-135-213.ap-northeast-1.compute.amazonaws.com:/tmp/slow_query_20231025021852.digest ./server1/digest
```
```
# upload server1
## mysql
```
scp -i "~/.ssh/id_rsa" -r ./server1/conf/mysqld.cnf isucon@ec2-3-112-135-213.ap-northeast-1.compute.amazonaws.com:/etc/mysql/mysql.conf.d/mysqld.cnf
```
## nginx
```
scp -i "~/.ssh/id_rsa" -r ./server1/conf/nginx.conf isucon@ec2-3-112-135-213.ap-northeast-1.compute.amazonaws.com:/etc/nginx/nginx.conf
```