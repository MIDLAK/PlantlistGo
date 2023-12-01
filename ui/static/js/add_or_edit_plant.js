
/* ДОБАВЛЕНИЕ РАСТЕНИЯ */

main()
plantData = main()
console.log(plantData.Promise)

async function main () {
	return plantData = await getplant()
}


//добавление публикации в список
const publicationAddButton = document.getElementById('publicationsButton')
// let publications = []
publicationAddButton.addEventListener('click', function () {
	//публикации (список)
	publication = document.getElementById('publications').value
	// publications.push(publication)
	plantData.publications.push(publication)
	elemNumber = Math.round(Math.random() * (1000000)) //не очень хорошо
	publicationsUL = document.getElementById('publicationsUL').insertAdjacentHTML('beforeend', '<li id="pubLi' + elemNumber + '"> <span id="pubSpan' + elemNumber  +'">'+ publication + '</span> <button type="button" class="btn btn-link" name="removePublicationsButtons" id="removePublication' + elemNumber + '">удалить</button> </li>')
	publicationsButtons()
})

//добавление меры сохранения вида
const saveMeasureButton = document.getElementById('saveMeasureButton')
// let measures = []
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

	// measures.push({"saveName": saveName,
	// 			   "description": saveDescription, 
	// 			   "startDate": startDate, 
	// 			   "endDate": endDate})
	plantData["saveMeasure"].push({"saveName": saveName,
								"description": saveDescription, 
								"startDate": startDate, 
								"endDate": endDate})

	elemNumber = Math.round(Math.random() * (1000000)) //не очень хорошо
	console.log(elemNumber)
	saveUL = document.getElementById('saveUL').insertAdjacentHTML('beforeend', '<li id="conservLi' + elemNumber + '"> <span id="conservSpan' + elemNumber + '">' + saveName + '</span> (с ' + startDate + ' по ' + endDate + ')' + '<button type="button" class="btn btn-link" name="removeConservation" id="removeConservationButton' + elemNumber + '">удалить</button>' + '</li>')
	measuresButtons()
})

//добавление места обитания
const placeButton = document.getElementById('placeButton')
//let places = []
placeButton.addEventListener('click', function () {
	// country = document.getElementById('country').value
	// city = document.getElementById('city').value
	// locality = document.getElementById('locality').value
	latitude = document.getElementById('gpsLatitude').value
	longitude = document.getElementById('gpsLongitude').value
	if (Math.abs(latitude) > 90 || Math.abs(longitude) > 180){
		alert('Ошибка в координатах (-90 <= широта <= 90 и -180 <= долгота <= 180)')
	} else {
		//places.push({'latitude': latitude, 'longitude': longitude})
		plantData["places"].push({'latitude': latitude, 'longitude': longitude})
		elemNumber = Math.round(Math.random() * (1000000)) //не очень хорошо
		placeUL = document.getElementById('placeUL').insertAdjacentHTML('beforeend', '<li> широта: <span id="spanLat'+ elemNumber  +'">'+ latitude  + '</span> долгота: <span id="spanLong' + elemNumber + '">' + longitude + '</span> <button type="button" class="btn btn-link" name="removePlacesButtons" id="removePlace'+ elemNumber +'">удалить</button> </li>')
		pointsButtons()
	}
})


measuresButtons()
function measuresButtons () {
	const removeSaveMeasuresButtons = document.getElementsByName('removeConservation')
	removeSaveMeasuresButtons.forEach(function (removeButton) {
		removeButton.addEventListener('click', function () {
			buttonNum = parseInt(removeButton.id.match(/\d+/))
			plantData["saveMeasure"].forEach(function (measure) {
				conservName = document.getElementById("conservSpan" + buttonNum).textContent
				if (measure.SaveName == conservName) {
					plantData["saveMeasure"].splice(plantData["saveMeasure"].indexOf(measure), 1) //удаление элемента
				}
			})

			//уборка отрисовки
			var elem = document.getElementById("conservLi" + buttonNum)
			elem.parentNode.removeChild(elem)
		})
	})
}

publicationsButtons()
function publicationsButtons () {
	const removePublicationsButtons = document.getElementsByName('removePublicationsButtons')
	removePublicationsButtons.forEach(function (pubButton) {
		pubButton.addEventListener('click', function () {
			buttonNum = parseInt(pubButton.id.match(/\d+/))
			plantData["publications"].forEach(function (publication) {
				pubName = document.getElementById("pubSpan" + buttonNum).textContent
				if (publication == pubName) {
					plantData["publications"].splice(plantData["publications"].indexOf(publication), 1) //удаление элемента
				}
			})

			//уборка отрисовки
			var elem = document.getElementById("pubLi" + buttonNum)
			elem.parentNode.removeChild(elem)
		})
	})
}

pointsButtons()
function pointsButtons () {
	const removePointsButtons = document.getElementsByName('removePlacesButtons')
	removePointsButtons.forEach(function (placeButton) {
		placeButton.addEventListener('click', function () {
			buttonNum = parseInt(placeButton.id.match(/\d+/))
			plantData["places"].forEach(function (place) {
				longitude = document.getElementById("spanLong" + buttonNum).textContent
				latitude = document.getElementById("spanLat" + buttonNum).textContent
				if (place.Latitude == latitude && place.Longitude == longitude) {
					plantData["places"].splice(plantData["places"].indexOf(place), 1) //удаление элемента
				}
			})

			//уборка отрисовки
			var elem = document.getElementById("placeLi" + buttonNum)
			elem.parentNode.removeChild(elem)
		})
	})
}

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
		publications: plantData["publications"],
		places: plantData["places"],
		saveMeasure: plantData["saveMeasure"],
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
			console.log(result)
		})
	}).catch((error) => {
		console.log(error)
	})
})

//POST-запрос на сервер
function getplant () {
	return new Promise((resolve, reject) => {
		plant = {
			plantName: document.getElementById('plantName').value
		}
		fetch('/getplant', {
			headers: {
				'Accept': 'application/json',
				'Content-Type': 'application/json'
			},
			method: 'POST',
			body: JSON.stringify(plant) //формирование JSON
		}).then((response) => {
			let dataResp
			response.text().then(function (respData) {
				result = JSON.parse(respData)
				console.log(result)
				resolve(result)
			})
		}).catch((error) => {
			console.log(error)
			return '{}'
		})
	});
}


