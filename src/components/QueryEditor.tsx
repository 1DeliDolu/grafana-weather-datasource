import React, { ChangeEvent } from 'react';
import { InlineField, Select, Stack } from '@grafana/ui';
import { QueryEditorProps } from '@grafana/data';
import { DataSource } from '../datasource';
import { MyDataSourceOptions, MyQuery ,stadte} from '../types';

type Props = QueryEditorProps<DataSource, MyQuery, MyDataSourceOptions>;

export function QueryEditor({ query, onChange, onRunQuery }: Props) {
  const onCityChange = (event: ChangeEvent<HTMLInputElement>) => {
    onChange({ ...query, city: event.target.value });
  };


  return (
    <Stack gap={0}>
      <InlineField label="City" labelWidth={16} tooltip={'Used'}>
        <Select
          width={16}
          onChange={city => onCityChange({ target: { value: city.value } } as ChangeEvent<HTMLInputElement>)}
          options={stadte}
          allowCustomValue
          isSearchable
        />
      </InlineField>
    </Stack>
  );
}
