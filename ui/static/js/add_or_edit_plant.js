
/* ДОБАВЛЕНИЕ РАСТЕНИЯ */

//добавление публикации в список
const publicationAddButton = document.getElementById('publicationsButton')
let publications = []
publicationAddButton.addEventListener('click', function () {
	//публикации (список)
	publication = document.getElementById('publications').value
	publications.push(publication)
	publicationsUL = document.getElementById('publicationsUL').insertAdjacentHTML('beforeend', '<li>'+ publication + '</li>')
})

//добавление меры сохранения вида
const saveMeasureButton = document.getElementById('saveMeasureButton')
let measures = []
saveMeasureButton.addEventListener('click', function () {
	saveName = document.getElementById('saveName').value
	saveDescription = document.getElementById('saveDescription').value
	startDate = document.getElementById('startDate').value
	if (startDate == '') {
		startDate = '...'
	}
	endDate = document.getElementById('endDate').value
	if (endDate == '') {
		endDate = '...'
	}

	measures.push({"saveName": saveName,
				   "description": saveDescription, 
				   "startDate": startDate, 
				   "endDate": endDate})

	saveUL = document.getElementById('saveUL').insertAdjacentHTML('beforeend', '<li>'+ saveName + ' (с ' + startDate + ' по ' + endDate + ')' + '</li>')
})

//добавление места обитания
const placeButton = document.getElementById('placeButton')
let places = []
placeButton.addEventListener('click', function () {
	// country = document.getElementById('country').value
	// city = document.getElementById('city').value
	// locality = document.getElementById('locality').value
	latitude = document.getElementById('gpsLatitude').value
	longitude = document.getElementById('gpsLongitude').value
	if (Math.abs(latitude) > 90 || Math.abs(longitude) > 180){
		alert('Ошибка в координатах (-90 <= широта <= 90 и -180 <= долгота <= 180)')
	} else {
		places.push({'latitude': latitude, 'longitude': longitude})
		placeUL = document.getElementById('placeUL').insertAdjacentHTML('beforeend', '<li> широта: '+ latitude  + ' долгота: ' + longitude + '</li>')
	}
})

//отравка JSON на сервер
const plantSubmitButton = document.getElementById('plantSubmitButton')
plantSubmitButton.addEventListener('click', function () {

	//фото
	// file = document.getElementById('customFile').files[0]
	// let reader = new FileReader()
	// reader.readAsDataURL(file)

	//называние растения
	plantName = document.getElementById('plantName').value
	plantLatinName = document.getElementById('plantLatinName').value

	//получение систематики
	sysElements = document.getElementsByName('systematization')
	domain = sysElements[0].value
	kingdom = sysElements[1].value
	department = sysElements[2].value
	plantClass = sysElements[3].value
	order = sysElements[4].value
	family = sysElements[5].value
	genus = sysElements[6].value

	// system = {"domain": domain,
	// 		  "kingdom": kingdom,
	// 		  "department": department,
	// 		  "class": plantClass,
	// 		  "order": order,
	// 		  "family": family,
	// 		  "genus": genus}

	//описание
	description = document.getElementById('plantDescription').value

	//статус
	plantStatus = document.getElementById('selectStatus').value

	//подготовка к отправке
	dataInput = {
		name: plantName,
		latinName: plantLatinName,
		domain: domain,
		kingdom: kingdom,
		department: department,
		class: plantClass,
		order: order,
		family: family,
		genus: genus,
		status: plantStatus,
		description: description,
		publications: publications,
		places: places,
		saveMeasure: measures,
	}

	console.log(JSON.stringify(dataInput))

	//POST-запрос на сервер
	fetch('/newplant', {
		headers: {
			'Accept': 'application/json',
			'Content-Type': 'application/json'
		},
		method: 'POST',
		body: JSON.stringify(dataInput) //формирование JSON
	}).then((response) => {
		let dataResp
		response.text().then(function (respData) {
			result = JSON.parse(respData)
			console.log(result)
		})
	}).catch((error) => {
		console.log(error)
	})
	
})

