<h3 align="center">Пример работы с Server-Sent Events (SSE) в приложении на React+Redux toolkit query</h3> 

<h4>Backend</h4>
****
Приложение принимает новые статусы групп через POST-эндпоинт **(/objects)**.
Полученные данные сохраняются и подготавливаются к обработке.
Новые статусы добавляются к уже существующему списку статусов.

**Server-Sent Events (SSE)**

Приложение предоставляет SSE-эндпоинт **(/events)** для стриминга обновлений клиентам.
Каждый статус обновляется и отправляется клиенту по одному с задержкой в 3 секунду между обновлениями (для большей наглядности).
Если клиент разрывает соединение во время стриминга, сервер завершает процесс отправки.

<h5>Согласованность данных</h5>
При обработке нескольких сущностей гарантируется согласованность данных, что означает, что данные останутся корректными даже в случае разрыва соединения.

**Эндпоинты**
<ul>
<li style="color: darksalmon">POST /objects: Принимает новые статусы групп.</li>
<li style="color: darksalmon">GET /events: Отправляет обновления статусов в режиме реального времени с помощью SSE.</li>
<li style="color: darksalmon">GET /connect: Возвращает текущие статусы всех групп.</li>
</ul>

<h4>Frontend</h4>
****
Приложение получает данные с эндпоинта **GET /connect** сохраняет их в кеш Redux.
Далее при отправке новых данных на **POST /objects** мы начинаем прослушивать сообщения 
от бекенда по **GET /events** и при получении нового сообщения обновляем данные в кеше.
В конце обработки мы с последним сообщением получаем флаг _processing_end_ после чего 
закрываем соединение. 
