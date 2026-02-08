## Summary

Closes #<!-- Issue番号 -->

<!-- 実験の目的を1-2行で記述 -->

## Environment

| Key | Value |
|:---|:---|
| Go | `go version` |
| OS/Arch | `GOOS/GOARCH` |
| CPU | |

## Results

### Size / Static Analysis

<!-- unsafe.Sizeof やオフセット等の静的検証結果 -->

| Struct | Size | Note |
|:---|---:|:---|
| A | -- B | |
| B | -- B | |

### Benchmark

```text
(go test -bench=. -benchmem -count=5 の出力を貼り付け)
```

### Summary Table

| Metric | Pattern A | Pattern B | Diff |
|:---|---:|---:|:---|
| | | | |

## Conclusion

<!-- verified / unexpected / inconclusive のいずれか -->

- **Result**: `result:___`
- <!-- 仮説に対する結論を1-2行で記述 -->
