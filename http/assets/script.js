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

	Chart.defaults.global.defaultFontColor = "rgba(21, 22, 23, 1)";

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

function timeSince(date) {
	var seconds = Math.floor((new Date() - date) / 1000);
	var interval = Math.floor(seconds / 31536000);
	if (interval > 1) {
		return interval + " years";
	}
	interval = Math.floor(seconds / 2592000);
	if (interval > 1) {
		return interval + " months";
	}
	interval = Math.floor(seconds / 86400);
	if (interval > 1) {
		return interval + " days";
	}
	interval = Math.floor(seconds / 3600);
	if (interval > 1) {
		return interval + " hours";
	}
	interval = Math.floor(seconds / 60);
	if (interval > 1 || interval == 0) {
		return interval + " minutes";
	}
	return interval + " minute";
}

var ws = new WebSocket("ws://" + window.location.host + "/socket");
var oldWaterLevelList = new(Array)
var oldVibrate = new(Boolean)
var oldUptime = new(String)
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
	if (oldVibrate != m.vibrate) {
		document.getElementsByClassName("vibrate")[0].textContent = m.vibrate;
		$(".vibrate").effect("bounce")
	}
	oldVibrate = m.vibrate;

	// Update running on text.
	// TODO: Do we really need to do this every message?
	document.getElementsByClassName("hostname")[0].textContent = m.hostname;


	// Update uptime text.
	uptime = timeSince(new Date(m.uptime))
	if (oldUptime != uptime) {
		document.getElementsByClassName("uptime")[0].textContent = uptime;
		$(".uptime").effect("bounce")
	}
	oldUptime = uptime;

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
