Деплой

1) **Нахожу список таргетов, куда деплоить**
Если команда запуска приложения server
    и атрибуты запроса tenant/project/application заданы
    и у прислонного ключа есть права на деплой, то
        в папке tenant/project/application
        смотрится файл app.yaml, если в нём задан объект app/repository, и обрабатываются все yaml из app/path или рута
Иначе (вариант для локального запуска)
    Прохожу по двум deployto папкам: в текущем и родительском каталоге.
    Собираю все yaml файлы из этих папок (в том числе подпапок), ищу в них:
    1.1) CRD envirement с подходящем именем, если не найден и ещё не на сервере, то команда на деплой на указанный envirement отправляется на сервер
    1.2) ищу в yaml файлах CRD target указанный в envirement (их может быть несколько)


2) **Поиск helm**
Формирую список папок, которые буду деплоить через  helm install
Для server - работаю в каталоге tenant/project/application, или в переопределённом в CRD app/repository app/path
папка helm - добавляю в список helms
ищу CRD "dependencies", они могут ссылаться на:
    2.1) папку содержащую Chart.yaml и values.yaml   - значит это уже helm, добавляю
    2.2) папку содержащую подпапку helm               - значит в ней helm, добавляю её
    2.2) папку содержащую подпапку deployto, тогда для неё рекурсивно вызывается Поиск helm(2). Одна папка может быть добавлена несколько раз, с разными 2.3) если alias совпадает, то деплоиться будет один раз (последний встретившийся)
    2.3) helm chart 

3) **Сборка образов, и их отправка в регистри**
прохожу по всем пап

4) непоследственно сам deploy
Для каждого target:
    получаю kubeconfig и namespace
    Для каждого компонента
        выполняю helm install для:
            1) helm указанный у target
            2) helm указанный у envirement
            3) helm'ы указанные в dependencies   -  один чарт может запускаться несколько раз, с разными values
            4) helm из текущего каталога
        и values:
            низкий приоритет ) values из текущего каталога
                            2) values указанный у envirement
            высокий приоритет) values указанный у target (они главные, пере)




--deployto  -  атрибут, позволяющий изменить рабочий каталог