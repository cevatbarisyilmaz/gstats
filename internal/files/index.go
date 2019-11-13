package files

var indexHTML = []byte(`<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>GStats</title>
		<link href="https://fonts.googleapis.com/css?family=Lato|Open+Sans&display=swap" rel="stylesheet">
		<link rel="stylesheet" href="style.css">
	</head>
	<body>
		<div id="main">
			<div id="current">
				<div class="current-element" id="current-unique-ips">
					<div class="current-header">Unique Visitors</div>
					<div id="current-unique-ips-value" class="current-value"></div>
				</div>
				<div class="current-element" id="current-connections">
					<div class="current-header">Connections</div>
					<div id="current-connections-value" class="current-value"></div>
				</div>
				<div class="current-element" id="current-inbound-width">
					<div class="current-header">Inbound Width</div>
					<div id="current-inbound-width-value" class="current-value"></div>
				</div>
				<div class="current-element" id="current-outbound-width">
					<div class="current-header">Outbound Width</div>
					<div id="current-outbound-width-value" class="current-value"></div>
				</div>
			</div>
			<div id="tab-selector">
				<div id="today" class="tab-selection selected">Today</div>
				<div id="yesterday" class="tab-selection">Yesterday</div>
				<div id="this-month" class="tab-selection">This Month</div>
				<div id="last-month" class="tab-selection">Last Month</div>
				<div id="this-year" class="tab-selection">This Year</div>
				<div id="last-year" class="tab-selection">Last Year</div>
			</div>
		</div>
		<script src="script.js"></script>
	</body>
</html>`)
