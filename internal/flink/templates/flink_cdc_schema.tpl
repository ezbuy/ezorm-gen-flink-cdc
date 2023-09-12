{{ define "flink_cdc_schema" -}}
{{- $from := .From }}
{{- $to := .To }}
CREATE TABLE {{ get $from.Args "connector" | snakecase }}_{{ $from.Database }}_{{ $from.Table }} (
	{{- range $i,$col := $from.Fields }}
	{{$col.GetName}} {{$col.GetType}},
	{{- end }}
	PRIMARY KEY ({{$from.PrimaryKey}}) NOT ENFORCED
) WITH (
	{{ $from.Args | formatArgs }}
);

CREATE TABLE {{ get $to.Args "connector" }}_{{ $to.Database }}_{{ $to.Table }} (
	{{- range $i,$col := $to.Fields }}
	{{$col.GetName}} {{$col.GetType}},
	{{- end }}
	PRIMARY KEY ({{$to.PrimaryKey}}) NOT ENFORCED
) WITH (
	{{ $to.Args | formatArgs }}
);

{{- end }}
