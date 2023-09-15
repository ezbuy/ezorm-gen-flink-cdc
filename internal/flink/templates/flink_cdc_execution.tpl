{{ define "flink_cdc_execution" -}}
EXECUTE STATEMENT SET
BEGIN
{{- range $exe := .}}
	INSERT INTO {{$exe.From.Connector | snakecase }}_{{ $exe.From.Database }}_{{$exe.From.Table}} SELECT * FROM {{$exe.To.Connector | snakecase }}_{{ $exe.To.Database }}_{{$exe.To.Table}};
{{- end }}
END
{{- end }}
