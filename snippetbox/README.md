### SNIPPETBOX
## Опис функціоналу
- Можливість створювати нові сніпети (збережені фрагменти тексту)
- Можливість ділитись створеними сніпетами
- Можливість переглядати створені сніпети
## Технології
- Реалізувати сервер з використанням Go
- Для створення бази даних використати MongoDB
## Використані бібліотеки
Для реалізації клієнт-серверної взаємодії використана бібліотека `net/http`.    
З використаннями даної бібліотеки формується servmux, що оброблює кожени варіант запиту клієнту та підбирає під нього відповідну handle функцію. Наприклад:

Для handle функції, що оброблює доступ до перегляду сніпету з конкретним id /snippet/view/{id}
~~~go
func snippetView(w http.ResponseWriter, r *http.Request) {
	// Extract the value of the id wildcard from the request using r.PathValue()
	// and try to convert it to an integer using the strconv.Atoi() function. If
	// it can't be converted to an integer, or the value is less than 1, we
	// return a 404 page not found response.
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
}
~~~
Буде наступна її реєстрація в servmux з подальшим запуском сервера
~~~go
// Use the http.NewServeMux() function to initialize a new servemux
mux := http.NewServeMux()

// Register handler function
mux.HandleFunc("GET /snippet/view/{id}", snippetView)

// Use the http.ListenAndServe() function to start a new web server. 
err := http.ListenAndServe("0.0.0.0:4000", mux)
log.Fatal(err)
~~~

## Структура проекту
Проекти містить 3 директорії:
- cmd - містить специфічний для програми код для виконуваних додатків. В нашому випадку це веб-додаток
- internal - містить допоміжний код, не специфічний для програми, який використовується в проекті. В нашому випадку використовується для зберігання потенційно багаторазово використовуваного коду, такого як допоміжні засоби перевірки та моделі баз даних Mongo
- ui - містить активи інтерфейсу користувача, які використовуються веб-додатком. Зокрема, директорія ui/html містить шаблони HTML, а директорія ui/static міститить статичні файли (CSS та зображення)