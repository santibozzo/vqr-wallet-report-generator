# VQR wallet report generator
## Generar binario
```
$ go mod tidy
$ go build
```

## Qué hace?
1. Toma reporte de administrador como input
2. Genera reporte de billetera (incluyendo solo registros APPROVED...)
3. Genera archivo con resultados:
    - Cantidad de registros en reporte de administrador
    - Cantidad de registros en reporte de billetera
    - Monto bruto total
    - Monto neto total
    - Comisión
    - IVA
    - Monto a pagar al administrador

## Manual de uso
1. En el mismo directorio del binario poner el reporte del administrador (tiene que terminar con `ARS_report.csv`)
2. Ejecutar binario
3. Se genera el reporte de la billetera (`<nombre>_wallet_report.csv`)
4. Se genera archivo de resultado (`<nombre>_wallet_report-RESULTS.txt`)