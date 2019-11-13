package files

var styleCSS = []byte(`body {
	margin: 0 0 5vw 0;
	padding: 0;
	font-family: 'Lato', sans-serif;
}
#main {
	width: 100%;
	margin: auto;
	display: flex;
	flex-direction: column;
	justify-content: center;
	align-items: center;
}
#current {
	max-width: 1080px;
	margin: auto;
	display: flex;
	flex-direction: row;
	justify-content: center;
	align-items: center;
	flex-wrap: wrap;
	width: 100%;
}
.current-element {
	min-width: calc(64px + 10%);
	padding: 8px;
	margin: 8px;
	flex-basis: auto;
	flex-grow: 1;
}
.current-header {
	font-size: 1.1em;
	padding: 4px;
}
.current-value {
	font-size: 1.25em;
	padding: 4px;
}

#tab-selector {
	width: 100%;
	display: flex;
	justify-content: center;
	align-items: center;
	font-size: 1em;
	background-color: black;
	padding: 2px;
	box-sizing: border-box;
}

.tab-selection {
	padding: 4px 8px;
	margin: 8px;
	min-width: calc(32px + 5%);
	text-align: center;
	background-color: rgb(32, 32, 32);
	border-radius: 8px;
	color: rgb(223, 223, 223);
	cursor: pointer;
}

.tab-selection:hover {
	background-color: rgb(16, 16, 16);
}

.tab-selection:active {
	background-color: rgb(8, 8, 8);
}

.selected {
	background-color: black;
}

.tab-container {
	margin: auto;
	max-width: 1080px;
	width: 100%;
}

table {
	border-collapse: collapse;
	width: 100%;
	font-size: 0.9em;
}

td {
	border: 2px solid white;
	text-align: left;
	padding: 4px;
}

th {
	border: 2px solid white;
	text-align: left;
	padding: 4px;
	font-weight: normal;
}

tr:nth-child(even) {
	background-color: rgb(241, 241, 241);
}

tr:nth-child(odd) {
	background-color: rgb(248, 248, 248);
}

.highlights-selector {
	display: flex;
	flex-direction: row;
	justify-content: center;
	align-items: center;
	font-size: 1.1em;
}

.highlights-selection {
	padding: 4px;
	margin: 4px;
	flex-grow: 1;
	max-width: calc(32px + 8%);
	text-align: center;
	cursor: pointer;
}

.highlights-selection:hover {
	background-color: rgb(242, 242, 242);
}

.highlights-selection:active {
	background-color: rgb(231, 231, 231);
}

.highlights-selection-selected {
	background-color: rgb(222, 222, 222);
}

`)
