## Создать JSON API для сайта объявлений
Для запуска проекта нужно выполнить команду: 
```bash
docker compose up
```
Так же есть возможность посмотреть содержимое PostgreSQL c помощью pgAdmin4. Для этого нужно перейти на [localhost:5050](local) с логином live@admin.com
 и паролем password
 
### Пример работы метода создания объявления:
```
Request:
    POST http://localhost:8080/advert/

    {
      "title" : "Cat",
      "description" : "The cat is a small carnivorous mammal. It is the only domesticated species in the family Felidae and often referred to as the domestic cat to distinguish it from wild members of the family. The cat is either a house cat, a farm cat or a feral cat; latter ranges freely and avoids human contact.",
      "price" : 124,
      "photo_link": ["https://icatcare.org/app/uploads/2018/07/Thinking-of-getting-a-cat.png", "https://img.huffingtonpost.com/asset/5dcc613f1f00009304dee539.jpeg?cache=QaTFuOj2IM&ops=crop_834_777_4651_2994%2Cscalefit_720_noupscale", "https://www.humanesociety.org/sites/default/files/styles/1660_max/public/2018/06/kittens-in-shelter-69469.jpg?itok=qa4Xlm8T"]
    }
```
```
Response:
{
    "ID": 1
}
```

### Метод получения конкретного объявления
```
Request:
    GET http://localhost:8080/advert/1?fields=description,photo_links
```
В запросе 1 -это ID объявления.
fields - это опциональное поле, которое принимает описание(description), ссылки на все фото(photo_links)
```
Response:
    {
        "title": "Cat",
        "description": "The cat is a small carnivorous mammal. It is the only domesticated species in the family Felidae and often referred to as the domestic cat to distinguish it from wild members of the family. The cat is either a house cat, a farm cat or a feral cat; latter ranges freely and avoids human contact.",
        "price": 124,
        "photo_link": [
            "https://icatcare.org/app/uploads/2018/07/Thinking-of-getting-a-cat.png",
            "https://img.huffingtonpost.com/asset/5dcc613f1f00009304dee539.jpeg?cache=QaTFuOj2IM&ops=crop_834_777_4651_2994%2Cscalefit_720_noupscale",
            "https://www.humanesociety.org/sites/default/files/styles/1660_max/public/2018/06/kittens-in-shelter-69469.jpg?itok=qa4Xlm8T"
        ]
    }
```
 
### Метод получения списка объявлений
```
Request:
    GET http://localhost:8080/adverts/?page=1&sort=time&sort_type=desc
```
В запросе параметры:

- page - это номер странице, на странице присутствует 10 объявлений

- sort - это поле принимает аргументы, по которым происходит сортировка, принимает price(цена) или time(дата создания)

- sort_type - это поле сортировки по возрастанию(ask), убыванию(DESC)
```
Response:
  [
      {
          "title": "Cat",
          "price": 124,
          "photo_link": [
              "https://icatcare.org/app/uploads/2018/07/Thinking-of-getting-a-cat.png"
          ]
      }
  ]
```


