# Teyit.link


teyit.link internet sayfalarının kopyasını alarak; çevirimiçi görüntüleyebileceğiniz sekilde kaydeder. 
Böylece sayfalar silinmiş olsa bile erişmeye devam edebilirsiniz.

teyit.link bir teyit.org aracıdır.

## Installation

### Requirements

You'll need to have a MySQL server running. At NOD, we use Docker containers in development. 
Here is a quick setup, if you already have Docker setup on your local machine.

```bash
docker run --name=teyitlinkdb -d -p 3306:3306 -e MYSQL_USER=link -e MYSQL_PASSWORD=root -e MYSQL_DATABASE=teyitlink mysql/mysql-server:5.7
```
