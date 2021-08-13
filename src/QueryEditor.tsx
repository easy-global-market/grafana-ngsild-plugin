import defaults from 'lodash/defaults';

import React, { ChangeEvent, PureComponent, MouseEvent } from 'react';
import { LegacyForms, Button, InlineFormLabel, Select } from '@grafana/ui';
import { QueryEditorProps, SelectableValue } from '@grafana/data';
import { DataSource } from './DataSource';
import { defaultQuery, MyDataSourceOptions, MyQuery, PanelQueryFormat } from './types';
import { getTemplateSrv } from '@grafana/runtime';
import { VariableModel } from '@grafana/data/types/templateVars';
interface QueryContext {
  query?: string;
  name: string;
}

const { FormField } = LegacyForms;

const FORMAT_OPTIONS: Array<SelectableValue<PanelQueryFormat>> = [
  { label: 'Table', value: PanelQueryFormat.Table },
  { label: 'World Map', value: PanelQueryFormat.WorldMap },
];
let isWorldMap = true;
let variables = (getTemplateSrv().getVariables() as unknown) as Array<VariableModel & QueryContext>;

type Props = QueryEditorProps<DataSource, MyQuery, MyDataSourceOptions>;
export class QueryEditor extends PureComponent<Props> {
  componentDidUpdate() {
    let currentVariables = getTemplateSrv().getVariables();
    if (currentVariables && currentVariables !== variables) {
      variables = currentVariables;
      this.isContextSet(currentVariables);
    }
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
    onChange({ ...query, entityType: event.target.value });
  };

  onValueFilterQueryChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onChange, query } = this.props;
    onChange({ ...query, valueFilterQuery: event.target.value });
  };

  onMetadataSelectorChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onChange, query } = this.props;
    onChange({ ...query, metadataSelector: event.target.value });
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

  //Check if a variable named 'context' exists
  isContextSet(currentVariables: QueryContext[]) {
    let found = false;
    currentVariables.forEach((variable: QueryContext) => {
      if (variable.name === 'context') {
        found = true;
        //set the content of 'context' dashboard variable to the context variable to get it in backend
        this.props.query.context = variable.query ? variable.query : '';
      }
    });
    if (!found) {
      this.props.query.context = '';
      throw new Error('Create a dashboard variable named "context" with your context');
    }
  }

  render() {
    const query = defaults(this.props.query, defaultQuery);
    const { entityId, entityType, valueFilterQuery, metadataSelector } = query;

    return (
      <div>
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
        <div className="gf-form-inline">
          <FormField
            labelWidth={11}
            inputWidth={20}
            value={metadataSelector || ''}
            onChange={this.onMetadataSelectorChange}
            label="Metadata Selector"
          />
        </div>
        <Button size="md" variant="secondary" onClick={this.onConfirm}>
          Confirm
        </Button>
      </div>
    );
  }
}
