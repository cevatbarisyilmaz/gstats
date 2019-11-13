package files

var scriptJS = []byte(`let main = document.getElementById("main");

let currentUniqueIPsValue = document.getElementById("current-unique-ips-value");
let currentConnections = document.getElementById("current-connections-value");
let currentInboundWidthValue = document.getElementById("current-inbound-width-value");
let currentOutboundWidthValue = document.getElementById("current-outbound-width-value");

let today = document.getElementById("today");
let yesterday = document.getElementById("yesterday");
let thisMonth = document.getElementById("this-month");
let lastMonth = document.getElementById("last-month");
let thisYear = document.getElementById("this-year");
let lastYear = document.getElementById("last-year");

setInterval(refreshCurrent, 4000);
refreshCurrent();
getData();

function refreshCurrent() {
	const url = document.location.pathname + "api/current";
	fetch(url).then(response => {
		return response.json();
	}).then(response => {
		currentUniqueIPsValue.innerText = response.UniqueIPs;
		currentConnections.innerText = response.Connections;
		currentInboundWidthValue.innerText = readableBytes(response.InboundBandwidth) + "/s";
		currentOutboundWidthValue.innerText = readableBytes(response.OutboundBandwidth) + "/s";
	});
}

function getData() {
	const url = document.location.pathname + "api/data";
	fetch(url).then(response => {
		return response.json();
	}).then(response => {
		let visibleDiv;
		let todayDiv = createTab(response.Today, "day");
		today.addEventListener("click", () => {
			visibleDiv.style.display = "none";
			todayDiv.style.display = "";
			visibleDiv = todayDiv;
			document.getElementsByClassName("selected")[0].classList.remove("selected");
			today.classList.add("selected");
		});
		visibleDiv = todayDiv;
		main.appendChild(todayDiv);
		let yesterdayDiv = createTab(response.Yesterday, "day");
		yesterdayDiv.style.display = "none";
		yesterday.addEventListener("click", () => {
			visibleDiv.style.display = "none";
			yesterdayDiv.style.display = "";
			visibleDiv = yesterdayDiv;
			document.getElementsByClassName("selected")[0].classList.remove("selected");
			yesterday.classList.add("selected");
		});
		main.appendChild(yesterdayDiv);
		let thisMonthDiv = createTab(response.ThisMonth, "month");
		thisMonthDiv.style.display = "none";
		thisMonth.addEventListener("click", () => {
			visibleDiv.style.display = "none";
			thisMonthDiv.style.display = "";
			visibleDiv = thisMonthDiv;
			document.getElementsByClassName("selected")[0].classList.remove("selected");
			thisMonth.classList.add("selected");
		});
		main.appendChild(thisMonthDiv);
		let lastMonthDiv = createTab(response.LastMonth, "month");
		lastMonthDiv.style.display = "none";
		lastMonth.addEventListener("click", () => {
			visibleDiv.style.display = "none";
			lastMonthDiv.style.display = "";
			visibleDiv = lastMonthDiv;
			document.getElementsByClassName("selected")[0].classList.remove("selected");
			lastMonth.classList.add("selected");
		});
		main.appendChild(lastMonthDiv);
		let thisYearDiv = createTab(response.ThisYear, "year");
		thisYearDiv.style.display = "none";
		thisYear.addEventListener("click", () => {
			visibleDiv.style.display = "none";
			thisYearDiv.style.display = "";
			visibleDiv = thisYearDiv;
			document.getElementsByClassName("selected")[0].classList.remove("selected");
			thisYear.classList.add("selected");
		});
		main.appendChild(thisYearDiv);
		let lastYearDiv = createTab(response.LastYear, "year");
		lastYearDiv.style.display = "none";
		lastYear.addEventListener("click", () => {
			visibleDiv.style.display = "none";
			lastYearDiv.style.display = "";
			visibleDiv = lastYearDiv;
			document.getElementsByClassName("selected")[0].classList.remove("selected");
			lastYear.classList.add("selected");
		});
		main.appendChild(lastYearDiv);
	});
}

function createTab(data, type) {
	let tabContainer = document.createElement("div");
	tabContainer.setAttribute("class", "tab-container");
	tabContainer.appendChild(createSubhighlights(data.SubHighlights, type));
	tabContainer.appendChild(createHighlights(data.Highlights));
	return tabContainer;
}

const months = [
	"January",
	"February",
	"March",
	"April",
	"May",
	"June",
	"July",
	"August",
	"September",
	"October",
	"November",
	"December",
];

function createHighlights(data) {
	let highlightsContainer = document.createElement("div");
	highlightsContainer.setAttribute("class", "highlights-container"); 
	
	let highlightsSelector = document.createElement("div");
	highlightsSelector.setAttribute("class", "highlights-selector");
	
	let selectedHeader;
	let visibleDiv;
	
	let requests = document.createElement("div");
	requests.textContent = "Requests";
	requests.setAttribute("class", "highlights-selection highlights-selection-selected");
	selectedHeader = requests;
	highlightsSelector.appendChild(requests);
	let referrers = document.createElement("div");
	referrers.textContent = "Referrers";
	referrers.setAttribute("class", "highlights-selection");
	highlightsSelector.appendChild(referrers);
	let countries = document.createElement("div");
	countries.textContent = "Countries";
	countries.setAttribute("class", "highlights-selection");
	highlightsSelector.appendChild(countries);
	let browsers = document.createElement("div");
	browsers.textContent = "Browsers";
	browsers.setAttribute("class", "highlights-selection");
	highlightsSelector.appendChild(browsers);
	let oss = document.createElement("div");
	oss.textContent = "OSs";
	oss.setAttribute("class", "highlights-selection");
	highlightsSelector.appendChild(oss);
	
	highlightsContainer.appendChild(highlightsSelector);
	
	let highlights = document.createElement("div");
	
	let requestsDiv = createRequestHighlights(data.RequestPaths);
	highlights.appendChild(requestsDiv);
	visibleDiv = requestsDiv;
	let referrersDiv = createCommonHighlights(data.Referrers, "Referrer");
	referrersDiv.style.display = "none";
	highlights.appendChild(referrersDiv);
	let countriesDiv = createCommonHighlights(data.Countries, "Country");
	countriesDiv.style.display = "none";
	highlights.appendChild(countriesDiv);
	let browsersDiv = createCommonHighlights(data.Browsers, "Browser");
	browsersDiv.style.display = "none";
	highlights.appendChild(browsersDiv);
	let ossDiv = createCommonHighlights(data.OSs, "OS");
	ossDiv.style.display = "none";
	highlights.appendChild(ossDiv);
	
	requests.addEventListener("click", () => {
		visibleDiv.style.display = "none";
		requestsDiv.style.display = "";
		visibleDiv = requestsDiv;
		selectedHeader.classList.remove("highlights-selection-selected");
		requests.classList.add("highlights-selection-selected");
		selectedHeader = requests;
	});
	referrers.addEventListener("click", () => {
		visibleDiv.style.display = "none";
		referrersDiv.style.display = "";
		visibleDiv = referrersDiv;
		selectedHeader.classList.remove("highlights-selection-selected");
		referrers.classList.add("highlights-selection-selected");
		selectedHeader = referrers;
	});
	countries.addEventListener("click", () => {
		visibleDiv.style.display = "none";
		countriesDiv.style.display = "";
		visibleDiv = countriesDiv;
		selectedHeader.classList.remove("highlights-selection-selected");
		countries.classList.add("highlights-selection-selected");
		selectedHeader = countries;
	});
	browsers.addEventListener("click", () => {
		visibleDiv.style.display = "none";
		browsersDiv.style.display = "";
		visibleDiv = browsersDiv;
		selectedHeader.classList.remove("highlights-selection-selected");
		browsers.classList.add("highlights-selection-selected");
		selectedHeader = browsers;
	});
	oss.addEventListener("click", () => {
		visibleDiv.style.display = "none";
		ossDiv.style.display = "";
		visibleDiv = ossDiv;
		selectedHeader.classList.remove("highlights-selection-selected");
		oss.classList.add("highlights-selection-selected");
		selectedHeader = oss;
	});
	
	highlightsContainer.appendChild(highlights);
	
	return highlightsContainer;
}

function createRequestHighlights(data) {
	let table = document.createElement("table");
	table.setAttribute("class", "highlight-container");
	
	let header = document.createElement("tr");
	
	let path = document.createElement("th");
	path.textContent = "Path";
	header.appendChild(path);
	
	let requests = document.createElement("th");
	requests.textContent = "Requests";
	header.appendChild(requests);
	
	let uniqueVisitors = document.createElement("th");
	uniqueVisitors.textContent = "Unique Visitors";
	header.appendChild(uniqueVisitors);
	
	let averageResponseTime = document.createElement("th");
	averageResponseTime.textContent = "Average Response Time";
	header.appendChild(averageResponseTime);
	
	let averageResponseSize = document.createElement("th");
	averageResponseSize.textContent = "Average Response Size";
	header.appendChild(averageResponseSize);
	
	let successfulStatusCodeRate = document.createElement("th");
	successfulStatusCodeRate.textContent = "Successful Status Code Rate";
	header.appendChild(successfulStatusCodeRate);
	
	let topNonsuccessfulStatusCode = document.createElement("th");
	topNonsuccessfulStatusCode.textContent = "Top Non-Successful Status Code";
	header.appendChild(topNonsuccessfulStatusCode);
	
	let topNonsuccessfulStatusCodeRate = document.createElement("th");
	topNonsuccessfulStatusCodeRate.textContent = "Top Non-Successful Status Code Rate";
	header.appendChild(topNonsuccessfulStatusCodeRate);
	
	table.appendChild(header);
	if(data == null) {
		return table;
	}
	for(let highlight of data) {
		let tr = document.createElement("tr");
		
		let path = document.createElement("td");
		path.innerText = highlight.Path;
		tr.appendChild(path);
		
		let requests = document.createElement("td");
		requests.innerText = highlight.Requests;
		tr.appendChild(requests);
		
		let uniqueVisitors = document.createElement("td");
		uniqueVisitors.innerText = highlight.UniqueIPs;
		tr.appendChild(uniqueVisitors);
		
		let averageResponseTime = document.createElement("td");
		averageResponseTime.innerText = readableTime(highlight.AverageResponseTime);
		tr.appendChild(averageResponseTime);
		
		let averageResponseSize = document.createElement("td");
		averageResponseSize.innerText = readableBytes(highlight.AverageResponseSize);
		tr.appendChild(averageResponseSize);
		
		let successfulStatusCodeRate = document.createElement("td");
		successfulStatusCodeRate.innerText =  (highlight.SuccessfulStatusCodeRate * 100).toFixed(2) + "%";
		tr.appendChild(successfulStatusCodeRate);
		
		let topNonSuccessfulStatusCode = document.createElement("td");
		topNonSuccessfulStatusCode.innerText = highlight.TopNonSuccessfulStatusCode == 0 ? "" : highlight.TopNonSuccessfulStatusCode;
		tr.appendChild(topNonSuccessfulStatusCode);
		
		let topNonSuccessfulStatusCodeRate = document.createElement("td");
		topNonSuccessfulStatusCodeRate.innerText =  highlight.TopNonSuccessfulStatusCode == 0 ? "" : (highlight.TopNonSuccessfulStatusCodeRate * 100).toFixed(2) + "%";
		tr.appendChild(topNonSuccessfulStatusCodeRate);
		
		table.appendChild(tr);
	}
	
	return table;
	
}

function createCommonHighlights(data, type) {
	let table = document.createElement("table");
	table.setAttribute("class", "highlight-container");
	
	let header = document.createElement("tr");
	
	let identifier = document.createElement("th");
	identifier.textContent = type;
	header.appendChild(identifier);
	
	let uniqueVisitors = document.createElement("th");
	uniqueVisitors.textContent = "Unique Visitors";
	header.appendChild(uniqueVisitors);
	
	let requests = document.createElement("th");
	requests.textContent = "Requests";
	header.appendChild(requests);
	
	table.appendChild(header);
	if(data == null) {
		return table;
	}
	for(let highlight of data) {
		let tr = document.createElement("tr");
		
		let identifier = document.createElement("td");
		identifier.innerText = highlight.Identifier == "" ? "Unknown" : (type === "Country" ? countries[highlight.Identifier] : highlight.Identifier);
		tr.appendChild(identifier);
		
		let uniqueVisitors = document.createElement("td");
		uniqueVisitors.innerText = highlight.UniqueIPs;
		tr.appendChild(uniqueVisitors);
		
		let requests = document.createElement("td");
		requests.innerText = highlight.Requests;
		tr.appendChild(requests);
		
		table.appendChild(tr);
	}
	
	return table;
}

function createSubhighlights(data, type) {
	let subhighlightContainer = document.createElement("table");
	subhighlightContainer.setAttribute("class", "subhighlight-container");
	let header = document.createElement("tr");
	let identifier = document.createElement("th");
	switch (type) {
		case "day":
			identifier.textContent = "Hours";
			break;
		case "month":
			identifier.textContent = "Days";
			break;
		case "year":
			identifier.textContent = "Months";
			break;
	}
	header.appendChild(identifier);
	let uniqueIPs = document.createElement("th");
	uniqueIPs.textContent = "Unique Visitors";
	header.appendChild(uniqueIPs);
	let requests = document.createElement("th");
	requests.textContent = "Requests";
	header.appendChild(requests);
	let inboundWidth = document.createElement("th");
	inboundWidth.textContent = "Inbound Width";
	header.appendChild(inboundWidth);
	let outboundWidth = document.createElement("th");
	outboundWidth.textContent = "Outbound Width";
	header.appendChild(outboundWidth);
	subhighlightContainer.appendChild(header);
	let index = 0;
	for(let highlight of data) {
		let subhighlight = document.createElement("tr");
		subhighlight.setAttribute("class", "subhighlight");
		let identifier =  document.createElement("td");
		identifier.setAttribute("class", "subhighlight-element");
		switch (type) {
			case "day":
				identifier.textContent = index.toString();
				break;
			case "month":
				identifier.textContent = (index + 1).toString();
				break;
			case "year":
				identifier.textContent = months[index];
				break;
		}
		subhighlight.appendChild(identifier);
		let uniqueIPs = document.createElement("td");
		uniqueIPs.setAttribute("class", "subhighlight-element");
		uniqueIPs.textContent = highlight.UniqueIPs;
		subhighlight.appendChild(uniqueIPs);
		let connections = document.createElement("td");
		connections.setAttribute("class", "subhighlight-element");
		connections.textContent = highlight.Requests;
		subhighlight.appendChild(connections);
		let inboundWidth = document.createElement("td");
		inboundWidth.setAttribute("class", "subhighlight-element");
		inboundWidth.textContent = readableBytes(highlight.InboundBandwidth);
		subhighlight.appendChild(inboundWidth);
		let outboundWidth = document.createElement("td");
		outboundWidth.setAttribute("class", "subhighlight-element");
		outboundWidth.textContent = readableBytes(highlight.OutboundBandwidth);
		subhighlight.appendChild(outboundWidth);
		subhighlightContainer.appendChild(subhighlight);
		index++
	}
	return subhighlightContainer;
}

function readableBytes(bytes) {
	if(bytes === 0) {
		return "0 B";
	}
	let i = Math.floor(Math.log(bytes) / Math.log(1024)),
	sizes = ["B", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"];
	return (bytes / Math.pow(1024, i)).toFixed(2) * 1 + " " + sizes[i];
}

function readableTime(t) {
	if(t < 1000) {
		return t + " ns";
	}
	t /= 1000;
	if(t < 1000) {
		return t.toFixed(2) + " μs";
	}
	t /= 1000;
	if(t < 1000) {
		return t.toFixed(2) + " ms";
	}
	t /= 1000;
	if(t < 60) {
		return t.toFixed(2) + " s";
	}
	t /= 60;
	if(t < 60) {
		return t.toFixed(1) + " m";
	}
	t /= 60;
	return t.toFixed(1) + " h";
}

let countries = {
    "AF": "Afghanistan",
    "AL": "Albania",
    "DZ": "Algeria",
    "AS": "American Samoa",
    "AD": "Andorra",
    "AO": "Angola",
    "AI": "Anguilla",
    "AQ": "Antarctica",
    "AG": "Antigua and Barbuda",
    "AR": "Argentina",
    "AM": "Armenia",
    "AW": "Aruba",
    "AU": "Australia",
    "AT": "Austria",
    "AZ": "Azerbaijan",
    "BS": "Bahamas",
    "BH": "Bahrain",
    "BD": "Bangladesh",
    "BB": "Barbados",
    "BY": "Belarus",
    "BE": "Belgium",
    "BZ": "Belize",
    "BJ": "Benin",
    "BM": "Bermuda",
    "BT": "Bhutan",
    "BO": "Bolivia",
    "BA": "Bosnia and Herzegovina",
    "BW": "Botswana",
    "BV": "Bouvet Island",
    "BR": "Brazil",
    "IO": "British Indian Ocean Territory",
    "BN": "Brunei Darussalam",
    "BG": "Bulgaria",
    "BF": "Burkina Faso",
    "BI": "Burundi",
    "KH": "Cambodia",
    "CM": "Cameroon",
    "CA": "Canada",
    "CV": "Cape Verde",
    "KY": "Cayman Islands",
    "CF": "Central African Republic",
    "TD": "Chad",
    "CL": "Chile",
    "CN": "China",
    "CX": "Christmas Island",
    "CC": "Cocos (Keeling) Islands",
    "CO": "Colombia",
    "KM": "Comoros",
    "CG": "Congo",
    "CD": "Congo, the Democratic Republic of the",
    "CK": "Cook Islands",
    "CR": "Costa Rica",
    "CI": "Cote D'Ivoire",
    "HR": "Croatia",
    "CU": "Cuba",
    "CY": "Cyprus",
    "CZ": "Czech Republic",
    "DK": "Denmark",
    "DJ": "Djibouti",
    "DM": "Dominica",
    "DO": "Dominican Republic",
    "EC": "Ecuador",
    "EG": "Egypt",
    "SV": "El Salvador",
    "GQ": "Equatorial Guinea",
    "ER": "Eritrea",
    "EE": "Estonia",
    "ET": "Ethiopia",
    "FK": "Falkland Islands (Malvinas)",
    "FO": "Faroe Islands",
    "FJ": "Fiji",
    "FI": "Finland",
    "FR": "France",
    "GF": "French Guiana",
    "PF": "French Polynesia",
    "TF": "French Southern Territories",
    "GA": "Gabon",
    "GM": "Gambia",
    "GE": "Georgia",
    "DE": "Germany",
    "GH": "Ghana",
    "GI": "Gibraltar",
    "GR": "Greece",
    "GL": "Greenland",
    "GD": "Grenada",
    "GP": "Guadeloupe",
    "GU": "Guam",
    "GT": "Guatemala",
    "GN": "Guinea",
    "GW": "Guinea-Bissau",
    "GY": "Guyana",
    "HT": "Haiti",
    "HM": "Heard Island and Mcdonald Islands",
    "VA": "Holy See (Vatican City State)",
    "HN": "Honduras",
    "HK": "Hong Kong",
    "HU": "Hungary",
    "IS": "Iceland",
    "IN": "India",
    "ID": "Indonesia",
    "IR": "Iran, Islamic Republic of",
    "IQ": "Iraq",
    "IE": "Ireland",
    "IL": "Israel",
    "IT": "Italy",
    "JM": "Jamaica",
    "JP": "Japan",
    "JO": "Jordan",
    "KZ": "Kazakhstan",
    "KE": "Kenya",
    "KI": "Kiribati",
    "KP": "North Korea",
    "KR": "South Korea",
    "KW": "Kuwait",
    "KG": "Kyrgyzstan",
    "LA": "Lao People's Democratic Republic",
    "LV": "Latvia",
    "LB": "Lebanon",
    "LS": "Lesotho",
    "LR": "Liberia",
    "LY": "Libya",
    "LI": "Liechtenstein",
    "LT": "Lithuania",
    "LU": "Luxembourg",
    "MO": "Macao",
    "MG": "Madagascar",
    "MW": "Malawi",
    "MY": "Malaysia",
    "MV": "Maldives",
    "ML": "Mali",
    "MT": "Malta",
    "MH": "Marshall Islands",
    "MQ": "Martinique",
    "MR": "Mauritania",
    "MU": "Mauritius",
    "YT": "Mayotte",
    "MX": "Mexico",
    "FM": "Micronesia, Federated States of",
    "MD": "Moldova, Republic of",
    "MC": "Monaco",
    "MN": "Mongolia",
    "MS": "Montserrat",
    "MA": "Morocco",
    "MZ": "Mozambique",
    "MM": "Myanmar",
    "NA": "Namibia",
    "NR": "Nauru",
    "NP": "Nepal",
    "NL": "Netherlands",
    "NC": "New Caledonia",
    "NZ": "New Zealand",
    "NI": "Nicaragua",
    "NE": "Niger",
    "NG": "Nigeria",
    "NU": "Niue",
    "NF": "Norfolk Island",
    "MK": "North Macedonia, Republic of",
    "MP": "Northern Mariana Islands",
    "NO": "Norway",
    "OM": "Oman",
    "PK": "Pakistan",
    "PW": "Palau",
    "PS": "Palestinian Territory, Occupied",
    "PA": "Panama",
    "PG": "Papua New Guinea",
    "PY": "Paraguay",
    "PE": "Peru",
    "PH": "Philippines",
    "PN": "Pitcairn",
    "PL": "Poland",
    "PT": "Portugal",
    "PR": "Puerto Rico",
    "QA": "Qatar",
    "RE": "Reunion",
    "RO": "Romania",
    "RU": "Russian Federation",
    "RW": "Rwanda",
    "SH": "Saint Helena",
    "KN": "Saint Kitts and Nevis",
    "LC": "Saint Lucia",
    "PM": "Saint Pierre and Miquelon",
    "VC": "Saint Vincent and the Grenadines",
    "WS": "Samoa",
    "SM": "San Marino",
    "ST": "Sao Tome and Principe",
    "SA": "Saudi Arabia",
    "SN": "Senegal",
    "SC": "Seychelles",
    "SL": "Sierra Leone",
    "SG": "Singapore",
    "SK": "Slovakia",
    "SI": "Slovenia",
    "SB": "Solomon Islands",
    "SO": "Somalia",
    "ZA": "South Africa",
    "GS": "South Georgia and the South Sandwich Islands",
    "ES": "Spain",
    "LK": "Sri Lanka",
    "SD": "Sudan",
    "SR": "Suriname",
    "SJ": "Svalbard and Jan Mayen",
    "SZ": "Swaziland",
    "SE": "Sweden",
    "CH": "Switzerland",
    "SY": "Syrian Arab Republic",
    "TW": "Taiwan",
    "TJ": "Tajikistan",
    "TZ": "Tanzania, United Republic of",
    "TH": "Thailand",
    "TL": "Timor-Leste",
    "TG": "Togo",
    "TK": "Tokelau",
    "TO": "Tonga",
    "TT": "Trinidad and Tobago",
    "TN": "Tunisia",
    "TR": "Turkey",
    "TM": "Turkmenistan",
    "TC": "Turks and Caicos Islands",
    "TV": "Tuvalu",
    "UG": "Uganda",
    "UA": "Ukraine",
    "AE": "United Arab Emirates",
    "GB": "United Kingdom",
    "US": "United States of America",
    "UM": "United States Minor Outlying Islands",
    "UY": "Uruguay",
    "UZ": "Uzbekistan",
    "VU": "Vanuatu",
    "VE": "Venezuela",
    "VN": "Viet Nam",
    "VG": "Virgin Islands, British",
    "VI": "Virgin Islands, U.S.",
    "WF": "Wallis and Futuna",
    "EH": "Western Sahara",
    "YE": "Yemen",
    "ZM": "Zambia",
    "ZW": "Zimbabwe",
    "AX": "Åland Islands",
    "BQ": "Bonaire, Sint Eustatius and Saba",
    "CW": "Curaçao",
    "GG": "Guernsey",
    "IM": "Isle of Man",
    "JE": "Jersey",
    "ME": "Montenegro",
    "BL": "Saint Barthélemy",
    "MF": "Saint Martin (French part)",
    "RS": "Serbia",
    "SX": "Sint Maarten (Dutch part)",
    "SS": "South Sudan",
    "XK": "Kosovo"
  };

`)
