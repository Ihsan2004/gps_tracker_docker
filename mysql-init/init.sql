-- MySQL'deki 'gps_user' kullanıcısına her yerden erişim izni verir
ALTER USER 'gps_user'@'%' IDENTIFIED BY 'gps_user_pass';

-- Kullanıcıya gps_tracker veritabanında tüm yetkileri ver
GRANT ALL PRIVILEGES ON gps_tracker.* TO 'gps_user'@'%';

-- Değişiklikleri etkinleştir
FLUSH PRIVILEGES;
