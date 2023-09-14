{{ define "flink_cdc_schema" -}}
{{- $from := .From }}
{{- $to := .To }}
CREATE TABLE {{ get $from.Args "connector" | snakecase }}_{{ $from.Database }}_{{ $from.Table }} (
	{{- range $i,$col := $from.Fields }}
	{{$col.GetName | backtick }} {{$col.GetType}},
	{{- end }}
	PRIMARY KEY ({{$from.PrimaryKey}}) NOT ENFORCED
) WITH (
	'database-name' = '{{.From.Database}}',
	'table-name' = '{{.From.Table}}',
	{{ $from.Args | formatArgs }}
);

CREATE TABLE {{ get $to.Args "connector" | snakecase }}_{{ $to.Database }}_{{ $to.Table }} (
	{{- range $i,$col := $to.Fields }}
	{{$col.GetName | backtick }} {{$col.GetType}},
	{{- end }}
	PRIMARY KEY ({{$to.PrimaryKey}}) NOT ENFORCED
) WITH (
	'database-name' = '{{.From.Database}}',
	'table-name' = '{{.From.Table}}',
	{{ $to.Args | formatArgs }}
);

INSERT INTO {{ get $to.Args "connector" }}_{{ $to.Database }}_{{ $to.Table }} SELECT * FROM {{ get $from.Args "connector" | snakecase }}_{{ $from.Database }}_{{ $from.Table }};

{{- end }}
