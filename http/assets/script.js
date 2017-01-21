var waterLevel = {
	labels: [],
	datasets: [{
		type: 'line',
		label: 'WATER LEVEL',
		backgroundColor: 'rgba(255, 99, 132, 0.5)',
		borderColor: 'rgba(255, 99, 132, 1)',
		data: [],
		pointRadius: 4,
		pointBackgroundColor: 'rgba(255, 255, 255, 1)',
		pointBorderWidth: 2,
	}]
};

window.onload = function() {
	var ctx = document.getElementById("canvas").getContext("2d");

	window.graph = new Chart(ctx, {
		type: 'line',
		data: waterLevel,
		options: {
			responsive: true,
			legend: {
				display: false
			},
			scales: {
				yAxes: [{
					gridLines: {
						color: 'rgba(243, 244, 246, 1)',
						drawBorder: false,
						zeroLineColor: 'rgba(243, 244, 246, 1)'
					}
				}],
				xAxes: [{
					gridLines: {
						color: 'rgba(243, 244, 246, 1)',
						drawBorder: false
					},
					ticks: {
						display: false
					}
				}]
			}
		}
	});
};

var ws = new WebSocket("ws://" + window.location.host + "/socket");
ws.onmessage = function(event) {
	ws.send(JSON.stringify({message: [0]}))

	var m = JSON.parse(event.data);

	// Update current div text.
	document.getElementById("waterLevel").textContent = m.message[m.message.length-1];

	// Update graph.
	waterLevel.labels = m.message
	waterLevel.datasets.forEach(function(dataset) {
		dataset.data = m.message;
	})
	window.graph.update();
}

