package html_report

var reportTemplate = `
<html>

<head>
	<meta charset="utf8"/>
</head>

<body>
	{{ range .}}
	<div>
		<h1>{{ .Name }}</h1>
		<table border="0">
			<tr align="left"> <th>Total</th> <td>{{ .Total }}</td></tr>
			<tr align="left"> <th>Slowest</th> <td>{{ .Slowest }}</td></tr>
			<tr align="left"> <th>Fastest</th> <td>{{ .Fastest }}</td></tr>
			<tr align="left"> <th>Average</th> <td>{{ .Average }}</td></tr>
			<tr align="left"> <th>Concurrent</th> <td>{{ .Concurrent }}</td></tr>
			<tr align="left"> <th>Requests</th> <td>{{ .RealRequest }}/{{ .Request }}</td></tr>
			<tr align="left"> <th>Requests/sec</th> <td>{{ .RealQPS }}/{{ .QPS }}</td></tr>
		</table>

		{{ if .ThroughputsImg }}
		<h2>Throughputs</h2>
		<p><img src="{{ .ThroughputsImg }}"/></p>
		{{ end }}

		{{ if .LatenciesImg }}
		<h2>Latencies</h2>
		<p><img src="{{ .LatenciesImg }}"/></p>
		{{ end }}

		{{ if .HistogramImg }}
		<h2>Histogram</h2>
		<p><img src="{{ .HistogramImg }}"/></p>
		{{ end }}
	
		{{ if .Latpcts }}
		<h2>Latency distribution</h2>
		{{ range .Latpcts }}
		<p>{{ .Percent }}% in {{ .Latency }}</p>
		{{ end }}
		{{ end }}
	
		{{ if .Errs }}
		<h2>Error distribution</h2>
		<table border="1" cellspacing="0" cellpadding="4">
			{{ range $key, $value := .Errs }}
			<tr align="left"><td>{{ $key }} </td> <td>{{ $value }}</td></tr>
			{{ end }}
		</table>
		{{ end }}
	</div>
	{{ end }}

</body>
</html>
`
