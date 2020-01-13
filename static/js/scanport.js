var table;
var tableInitialized = false;
$(document).ready(function() {
	if (!tableInitialized) {
		tableInitialized = true;
		table = $('#scanport').DataTable({
			columns: [
				{ data: 'Hostname' },
				{ data: 'IPAddress' },
				{ data: 'Protocol' },
				{ data: 'Port' },
				{ data: 'Status' }
			],
			paging: false,
			createdRow: (row, data, dataIndex, cells) => {
				if (data['Status'] == `open`) {
					$(cells[4]).addClass('openPort');
				} else {
					$(cells[4]).addClass('closePort');
				}
			}
		});
	}
	$('.testButton').click(function() {
		scanports();
	});
});

function scanports() {
	var hostname = $('#hostname').val();
	var tcp_port = $('#tcp_port').val();
	var udp_port = $('#udp_port').val();
	var data;

	if (!hostname || (!tcp_port && !udp_port)) {
		console.log('You need a hostname and at least one protocol selected.');
		window.alert('You need a hostname and at least one protcol selected.');
		return;
	} else if (!tcp_port) {
		console.log('UDP only');
		data = `hostname=${hostname}&udp=${udp_port}`;
	} else if (!udp_port) {
		data = `hostname=${hostname}&tcp=${tcp_port}`;
	} else {
		data = `hostname=${hostname}&tcp=${tcp_port}&udp=${udp_port}`;
	}
	$('.testButton').html('<div class="ld ld-ring ld-spin"></div>');
	$.ajax({
		url: '/',
		type: 'POST',
		data: `hostname=${hostname}&tcp=${tcp_port}&udp=${udp_port}`,
		dataType: 'json',
		success: function(response) {
			// console.log(response);
			// console.log(response[0]);
			table.clear().rows.add(response).draw(); //clear results of table then redraw
			$('.testButton').html('TEST');
		},
		error: function(xhr, status, error) {
			$('.testButton').html('TEST');
			window.alert(`💩 Oops! I dah poopeded. :() (${error}).`);
		}
	});
}
