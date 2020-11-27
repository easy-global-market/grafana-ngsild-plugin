import { DataQuery, DataSourceJsonData } from '@grafana/data';

export interface MyQuery extends DataQuery {
  queryText?: string;
  format?: string;
  attribute?: string;
}

export const defaultQuery: Partial<MyQuery> = {};

/**
 * These are options configured for each DataSource instance
 */
export interface MyDataSourceOptions extends DataSourceJsonData {
  authServerUrl?: string;
  resource?: string;
  clientId?: string;
  contextBrokerUrl?: string;
}

/**
 * Value that is used in the backend, but never sent over HTTP to the frontend
 */
export interface MySecureJsonData {
  clientSecret?: string;
}

export enum PanelQueryFormat {
  Table = 'table',
  WorldMap = 'worldmap',
}
