---
title: "Модуль linstor: настройки"
force_searchable: true
---

{% include module-bundle.liquid %}

> Работоспособность модуля гарантируется только для стоковых ядер, поставляемых вместе с дистрибутивами перечисленными в [списке поддерживаемых ОС](../../supported_versions.html#linux).
> Работоспособность модуля на других ядрах возможна, но не гарантируется.

После включения модуля кластер автоматически настраивается на использование LINSTOR и остается только сконфигурировать хранилище.

Сам модуль не требует настройки и не имеет параметров. Однако отдельные его функции могут потребовать указания мастер-пароля.  
Для того чтобы задать мастер-пароль, создайте Secret в пространстве имен `d8-system`:

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: linstor-passphrase
  namespace: d8-system
immutable: true
stringData:
  MASTER_PASSPHRASE: *!пароль* # Мастер-пароль для LINSTOR
```

> **Внимание!** Подойдите ответственно к выбору мастер-пароля для LINSTOR. При его потере зашифрованные данные окажутся недоступными.

## Конфигурация хранилища LINSTOR

Конфигурация LINSTOR в Deckhouse осуществляется посредством назначения специального тега `linstor-<имя_пула>` на LVM-группу томов или LVMThin-пул.  

1. Выберите имя тега.

   Имя тега должно быть уникальным в пределах одного узла. Поэтому каждый раз, прежде чем назначить новый тег, убедитесь в отсутствии этого тега у других групп томов и пулов.

   Выполните следующие команды, чтобы вывести список групп томов и пулов:

   ```shell
   # LVM пулы
   vgs -o+tags | awk 'NR==1;$NF~/linstor-/'
   # LVMThin пулы
   lvs -o+tags | awk 'NR==1;$NF~/linstor-/'
   ```

1. Добавьте пулы.

   Создайте пулы хранения на всех узлах, где вы планируете хранить ваши данные. Используйте одинаковые имена пулов хранения на разных узлах, если хотите иметь для них общий StorageClass.

   - Чтобы добавить пул **LVM** создайте группу томов с тегом `linstor-<имя_пула>`, либо добавьте тег `linstor-<имя_пула>` существующей группе.

     Пример команды создания группы томов `vg0` с тегом `linstor-data`:

     ```shell
     vgcreate vg0 /dev/nvme0n1 /dev/nvme1n1 --add-tag linstor-data
     ```

     Пример команды добавления тега `linstor-data` существующей группе томов `vg0`:

     ```shell
     vgchange vg0 --add-tag linstor-data
     ```

   - Чтобы добавить пул **LVMThin** создайте LVMThin-пул с тегом `linstor-<имя_пула>`.

     Пример команды создания LVMThin-пула `vg0/thindata` с тегом `linstor-data`:

     ```shell
     vgcreate vg0 /dev/nvme0n1 /dev/nvme1n1
     lvcreate -l 100%FREE -T vg0/thindata --add-tag linstor-thindata
     ```

     > Обратите внимание, что сама группа томов не обязана содержать какой-либо тег.

1. Проверьте создание StorageClass.

   Когда все пулы хранения будут созданы, появятся три новых StorageClass'а. Проверьте что они создались, выполнив в кластере Kubernetes команду:

   ```shell
   kubectl get storageclass
   ```

   Пример вывода списка StorageClass:

   ```shell
   $ kubectl get storageclass
   NAME                   PROVISIONER                  AGE
   linstor-data-r1        linstor.csi.linbit.com       143s
   linstor-data-r2        linstor.csi.linbit.com       142s
   linstor-data-r3        linstor.csi.linbit.com       142s
   ```

   Каждый StorageClass можно использовать для создания томов соответственно с одной, двумя или тремя репликами в ваших пулах хранения.

При необходимости изучите пример [расширенной конфигурации LINSTOR](advanced_usage.html), но мы рекомендуем придерживаться приведенного выше упрощённого руководства.
