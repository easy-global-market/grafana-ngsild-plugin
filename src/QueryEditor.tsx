import defaults from 'lodash/defaults';

import React, { ChangeEvent, PureComponent, MouseEvent } from 'react';
import { LegacyForms, Button, InlineFormLabel, Select } from '@grafana/ui';
import { QueryEditorProps, SelectableValue } from '@grafana/data';
import { DataSource } from './DataSource';
import { defaultQuery, MyDataSourceOptions, MyQuery, PanelQueryFormat } from './types';

const { FormField } = LegacyForms;

const FORMAT_OPTIONS: Array<SelectableValue<PanelQueryFormat>> = [
  { label: 'Table', value: PanelQueryFormat.Table },
  { label: 'World Map', value: PanelQueryFormat.WorldMap },
];
var isWorldMap = true;

type Props = QueryEditorProps<DataSource, MyQuery, MyDataSourceOptions>;

export class QueryEditor extends PureComponent<Props> {
  onQueryTextChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onChange, query } = this.props;
    onChange({ ...query, queryText: event.target.value });
    //onRunQuery();
  };

  onAttributeChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onChange, query } = this.props;
    onChange({ ...query, attribute: event.target.value });
  };

  onConfirm = (event: MouseEvent) => {
    const { onRunQuery } = this.props;
    onRunQuery();
  };

  getFormatOption = () => {
    return FORMAT_OPTIONS.find(v => v.value === this.props.query.format);
  };

  onFormatChange = (option: SelectableValue<PanelQueryFormat>) => {
    const { query, onChange } = this.props;
    if (option.value) {
      onChange({ ...query, format: option.value });
      isWorldMap = option.value === 'worldmap';
    }
  };

  render() {
    const query = defaults(this.props.query, defaultQuery);
    const { queryText } = query;

    return (
      <>
        <div className="gf-form-inline">
          <FormField
            inputWidth={25}
            labelWidth={6}
            value={queryText || ''}
            onChange={this.onQueryTextChange}
            label="entityId"
            placeholder="urn:ngsi-ld: ..."
          />
          <InlineFormLabel width={6}>Format</InlineFormLabel>
          <Select
            isSearchable={false}
            width={20}
            options={FORMAT_OPTIONS}
            onChange={this.onFormatChange}
            value={this.getFormatOption()}
          />

          <Button size="md" variant="secondary" onClick={this.onConfirm}>
            Confirm
          </Button>
        </div>
        {isWorldMap && (
          <FormField
            labelWidth={11}
            inputWidth={20}
            label="Attribute to use as a metric"
            value={query.attribute || ''}
            onChange={this.onAttributeChange}
          />
        )}
      </>
    );
  }
}
