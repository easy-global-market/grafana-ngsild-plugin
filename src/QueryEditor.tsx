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
  //Init context with "$context" to replace it with dashboard variable
  componentDidMount() {
    this.onContextChange;
  }

  onEntityIdChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onChange, query } = this.props;
    onChange({ ...query, entityId: event.target.value });
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

  onEntityTypeChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onChange, query } = this.props;
    onChange({ ...query, entityType: event.target.value, context: '$context' });
  };

  onValueFilterQueryChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onChange, query } = this.props;
    onChange({ ...query, valueFilterQuery: event.target.value });
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

  onContextChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onChange, query } = this.props;
    onChange({ ...query, context: '$context' });
  };

  render() {
    const query = defaults(this.props.query, defaultQuery);
    const { entityId, entityType, valueFilterQuery } = query;

    return (
      <>
        <div className="gf-form-inline">
          <FormField
            labelWidth={11}
            inputWidth={20}
            value={entityId || ''}
            onChange={this.onEntityIdChange}
            label="Entity Identifier"
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
        <div className="gf-form-inline">
          <FormField
            labelWidth={11}
            inputWidth={20}
            value={entityType || ''}
            onChange={this.onEntityTypeChange}
            label="Entity Type"
          />
          <FormField
            labelWidth={11}
            inputWidth={20}
            value={valueFilterQuery || ''}
            onChange={this.onValueFilterQueryChange}
            tooltip="An expression conform to the NGSI-LD query language"
            placeholder="minValue>1;maxValue=5"
            label="Value Filter Query"
          />
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
