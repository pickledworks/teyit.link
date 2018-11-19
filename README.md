### WORK IN PROGRESS
This project currently works fine but should be used with caution since it's fairly new.

Also, there is no multi-language feature yet and this project's primary language is Turkish.

# teyit.link archiver

teyit.link is a web applicatio and a AWS Lambda function to archive websites. 
This project is build for Teyit.org (fact checking organization from Turkey) by NOD Digital.

**Turkish Description:**

teyit.link internet sayfalarının kopyasını alarak; çevirimiçi görüntüleyebileceğiniz sekilde kaydeder. 
Böylece sayfalar silinmiş olsa bile erişmeye devam edebilirsiniz.

teyit.link bir teyit.org aracıdır.

## Quick Installation

If you want to try teyit.link, we ship a Dockerfile that we actually use in production. 

### Requirements

You'll need to have a MySQL server running. At NOD, we use Docker containers in development. 
Here is a quick setup, if you already have Docker setup on your local machine.

```bash
docker run --name=teyitlinkdb -d -p 3306:3306 -e MYSQL_USER=link -e MYSQL_PASSWORD=root -e MYSQL_DATABASE=teyitlink mysql/mysql-server:5.7 --character-set-server=utf8mb4 --collation-server=utf8mb4_unicode_ci
```

## License

This project is licensed with MIT. 
