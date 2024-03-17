Перед запуском необходимо поднять [сервер](https://github.com/realPointer/YouTube-thumbnails-downloader)

По умолчанию запросы идут через порт `50052`

Скомпилировать бинарный файл и при запуске передать ссылки на YouTube видео в формате https://www.youtube.com/watch?v=ID и до ?v=ID
```zsh
# Компиляция
go build -o downloader main.go

# Запуск
./downloader <LINKS>
```
Флаги:

`--output` - место куда сохранить обложки (необходимо предвратительно создать)

`--async` - выполнить загрузку асинхронно

Примеры выполнения:
![image](https://github.com/realPointer/thumbnails-downloader-cli/assets/50529632/1c279683-cba0-44df-a9e4-17043cb14d59)

![image](https://github.com/realPointer/thumbnails-downloader-cli/assets/50529632/5436c9d7-748f-4e73-8ad3-f82e532e5701)

![image](https://github.com/realPointer/thumbnails-downloader-cli/assets/50529632/2c467b1d-1388-43f9-a394-610aaa3016c3)


