import React, { ChangeEvent, PureComponent } from 'react';
import { LegacyForms } from '@grafana/ui';
import { DataSourcePluginOptionsEditorProps } from '@grafana/data';
import { MyDataSourceOptions, MySecureJsonData } from './types';

const { SecretFormField, FormField } = LegacyForms;

interface Props extends DataSourcePluginOptionsEditorProps<MyDataSourceOptions> {}

interface State {}

export class ConfigEditor extends PureComponent<Props, State> {
  onTokenUrlChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onOptionsChange, options } = this.props;
    const jsonData = {
      ...options.jsonData,
      tokenUrl: event.target.value,
    };
    onOptionsChange({ ...options, jsonData });
  };

  onResourceChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onOptionsChange, options } = this.props;
    const jsonData = {
      ...options.jsonData,
      resource: event.target.value,
    };
    onOptionsChange({ ...options, jsonData });
  };

  onClientIdChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onOptionsChange, options } = this.props;
    const jsonData = {
      ...options.jsonData,
      client_id: event.target.value,
    };
    onOptionsChange({ ...options, jsonData });
  };

  onApiUrlChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onOptionsChange, options } = this.props;
    const jsonData = {
      ...options.jsonData,
      apiUrl: event.target.value,
    };
    onOptionsChange({ ...options, jsonData });
  };

  // Secure field (only sent to the backend)
  onClientSecretChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onOptionsChange, options } = this.props;
    onOptionsChange({
      ...options,
      secureJsonData: {
        client_secret: event.target.value,
      },
    });
  };

  onResetClientSecretKey = () => {
    const { onOptionsChange, options } = this.props;
    onOptionsChange({
      ...options,
      secureJsonFields: {
        ...options.secureJsonFields,
        client_secret: false,
      },
      secureJsonData: {
        ...options.secureJsonData,
        client_secret: '',
      },
    });
  };

  render() {
    const { options } = this.props;
    const { jsonData, secureJsonFields } = options;
    const secureJsonData = (options.secureJsonData || {}) as MySecureJsonData;

    return (
      <div className="gf-form-group">
        <div className="gf-form">
          <FormField
            label="SSO URL"
            labelWidth={8}
            inputWidth={22}
            onChange={this.onTokenUrlChange}
            value={jsonData.tokenUrl || ''}
            placeholder="https://my.sso.org"
          />
        </div>

        <div className="gf-form">
          <FormField
            label="Token endpoint path"
            labelWidth={8}
            inputWidth={22}
            onChange={this.onResourceChange}
            value={jsonData.resource || ''}
            placeholder="/path/to/token/endpoint"
          />
        </div>

        <div className="gf-form">
          <FormField
            label="Client id"
            labelWidth={8}
            inputWidth={22}
            onChange={this.onClientIdChange}
            value={jsonData.client_id || ''}
            tooltip="OAuth2 client id to be used by the plugin"
          />
        </div>

        <div className="gf-form-inline">
          <div className="gf-form">
            <SecretFormField
              isConfigured={(secureJsonFields && secureJsonFields.client_secret) as boolean}
              value={secureJsonData.client_secret || ''}
              label="Client secret"
              tooltip="OAuth client secret to be used by the plugin"
              placeholder=""
              labelWidth={8}
              inputWidth={22}
              onReset={this.onResetClientSecretKey}
              onChange={this.onClientSecretChange}
            />
          </div>
        </div>

        <div className="gf-form">
          <FormField
            label="NGSI-LD API URL"
            labelWidth={8}
            inputWidth={22}
            onChange={this.onApiUrlChange}
            value={jsonData.apiUrl || ''}
            placeholder="https://my.context-brocker.org"
          />
        </div>
      </div>
    );
  }
}
