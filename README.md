Приложение работает на порту 8080. 
Запускается через docker-compose up.
Программа имеет следующие энпоинты.
GET /api/ping: Проверка состояния.
GET /api/tenders: Получение спика тендеров с возможностью фильтрации по типу услуг (service_type), offset и limit, через query.
POST /api/tender/new: Создание нового тендера. Создавать могут только пользователи от имени своей организации.
GET /api/tenders/my: Получение спика тендеров пользователя с offset и limit, через query. Получать список могут только пользователи, указывается username через query.
GET /api/tenders/{tenderId}/status: Получение текущего статуса тендера. Status указывается через query. Получать список могут только пользователи, указывается username через query. 
PUT /api/tenders/{tenderId}/status: Изменения статуса тендера. Status указывается через query. Получать список могут только пользователи, указывается username через query. 
PATCH /api/tenders/{tenderId}/edit: Редактирование тендера. Можно редактировать такие параметр как: name, description, serviceType. Получать список могут только пользователи, указывается username через query. 
PUT /api/tenders/{tenderId}/rollback/{version}: Откат версии тендера к указанной версии. Получать список могут только пользователи, указывается username через query. 