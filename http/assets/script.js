var waterLevel = {
	labels: [],
	datasets: [{
		type: "line",
		label: "Water level",
		backgroundColor: "rgba(255, 99, 132, 0.5)",
		borderColor: "rgba(255, 99, 132, 1)",
		data: [],
		pointRadius: 4,
		pointBackgroundColor: "rgba(255, 255, 255, 1)",
		pointBorderWidth: 2,
	}]
};

window.onload = function() {
	var ctx = document.getElementById("canvas").getContext("2d");

	window.graph = new Chart(ctx, {
		type: "line",
		data: waterLevel,
		options: {
			responsive: true,
			legend: {
				display: false
			},
			scales: {
				yAxes: [{
					gridLines: {
						color: "rgba(243, 244, 246, 1)",
						drawBorder: false,
						zeroLineColor: "rgba(243, 244, 246, 1)"
					},
					ticks: {
						min: 0,
						max: 100
					}
				}],
				xAxes: [{
					gridLines: {
						color: "rgba(243, 244, 246, 1)",
						drawBorder: false,
						zeroLineColor: "rgba(243, 244, 246, 1)"
					},
					ticks: {
						display: false,
					}
				}]
			},
			tooltips: {
				backgroundColor: "rgba(53, 55, 58, 1)",
				displayColors: false,
			}
		}
	});
};

var ws = new WebSocket("ws://" + window.location.host + "/socket");
var oldWaterLevelList = new(Array)
var oldVibrate = new(Boolean)
ws.onmessage = function(event) {
	var m = JSON.parse(event.data);

	// Update water level text.
	l = m.waterLevelList.length
	if (oldWaterLevelList != m.waterLevelList[l-1]) {
		document.getElementsByClassName("waterLevel")[0].textContent = m.waterLevelList[l-1];
		$(".waterLevel").effect("bounce")
	}
	oldWaterLevelList = m.waterLevelList[l-1];

	// Update vibrate text.
	console.debug(oldVibrate)
	if (oldVibrate != m.vibrate) {
		document.getElementsByClassName("vibrate")[0].textContent = m.vibrate;
		$(".vibrate").effect("bounce")
	}
	oldVibrate = m.vibrate;

	// Update running on text.
	document.getElementsByClassName("ip")[0].textContent = m.ip.substring(0, 7)+"...";

	// Update vibrate text.
	console.debug(oldVibrate)
	if (oldVibrate != m.vibrate) {
		document.getElementsByClassName("vibrate")[0].textContent = m.vibrate;
		$(".vibrate").effect("bounce")
	}
	oldVibrate = m.vibrate;
	// Update graph.
	var labels = [];
	for (var i = 0; i != l; i++) labels.push("")
	waterLevel.labels = labels
	waterLevel.datasets.forEach(function(dataset) {
		dataset.data = m.waterLevelList;
	})
	window.graph.update();

	// Send a message to the server (that way the server can see if
	// it needs to keep the connection open.
	ws.send(JSON.stringify({vibrate: true}))
}
