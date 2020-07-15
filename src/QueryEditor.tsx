import defaults from 'lodash/defaults';

import React, { ChangeEvent, PureComponent, MouseEvent } from 'react';
import { LegacyForms, Button } from '@grafana/ui';
import { QueryEditorProps } from '@grafana/data';
import { DataSource } from './DataSource';
import { defaultQuery, MyDataSourceOptions, MyQuery } from './types';

const { FormField } = LegacyForms;

type Props = QueryEditorProps<DataSource, MyQuery, MyDataSourceOptions>;

export class QueryEditor extends PureComponent<Props> {
  onQueryTextChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onChange, query } = this.props;
    onChange({ ...query, queryText: event.target.value });
    //onRunQuery();
  };

  onConfirm = (event: MouseEvent) => {
    const { onRunQuery } = this.props;
    onRunQuery();
  };

  render() {
    const query = defaults(this.props.query, defaultQuery);
    const { queryText } = query;

    return (
      <div className="gf-form">
        <FormField
          inputWidth={25}
          labelWidth={6}
          value={queryText || ''}
          onChange={this.onQueryTextChange}
          label="entityId"
          tooltip="urn:ngsi-ld: ..."
        />

        <Button size="md" variant="secondary" onClick={this.onConfirm}>
          Confirm
        </Button>
      </div>
    );
  }
}
