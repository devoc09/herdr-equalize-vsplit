# Equalize VSplit

現在のペインを右方向に分割し、同列のペイン幅を均等にする Herdr プラグインです。

## 必要環境

- [Herdr](https://herdr.dev) v0.7.0 以上（macOS）
- [Go](https://go.dev) 1.26.2 以上（ビルド時）

## インストール

```console
herdr plugin install devoc09/herdr-equalize-vsplit
```

インストールが完了すると、プラグイン ID `equalize-vsplit` で登録されます。

## 使い方

### アクションの確認

```console
herdr plugin action list --plugin equalize-vsplit
```

### アクションの実行

```console
herdr plugin action invoke split --plugin equalize-vsplit
```

1 回の実行でペインが右に分割され、左右 2 列が均等幅になります。
分割された右側のペインで再度実行すると、3 列が均等幅になります。

### キーバインドの割り当て（任意）

`~/.config/herdr/config.toml` に以下を追加すると、ショートカットキーからアクションを実行できます。

```toml
[[keys.command]]
key = "prefix+\\"
type = "plugin_action"
command = "equalize-vsplit.split"
description = "Split and equalize columns"
```

## アンインストール

```console
herdr plugin uninstall equalize-vsplit
```

またはインストール時のソース指定でも解除できます。

```console
herdr plugin uninstall devoc09/herdr-equalize-vsplit
```

## 開発

### ローカルリンク

```console
herdr plugin link /path/to/herdr-equalize-vsplit
```

### テスト

```console
go test ./...
```

### ビルド

```console
go build -o bin/herdr-equalize-vsplit .
```
