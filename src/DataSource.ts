import { DataSourceInstanceSettings } from '@grafana/data';
import { DataSourceWithBackend } from '@grafana/runtime';
import { MyDataSourceOptions, MyQuery } from './types';
import { getTemplateSrv } from '@grafana/runtime';

export class DataSource extends DataSourceWithBackend<MyQuery, MyDataSourceOptions> {
  constructor(instanceSettings: DataSourceInstanceSettings<MyDataSourceOptions>) {
    super(instanceSettings);
  }
  applyTemplateVariables(query: MyQuery) {
    const templateSrv = getTemplateSrv();
    return {
      ...query,
      entityId: query.entityId ? templateSrv.replace(query.entityId) : '',
      attribute: query.attribute ? templateSrv.replace(query.attribute) : '',
      context: query.context ? templateSrv.replace(query.context) : '',
    };
  }
}
