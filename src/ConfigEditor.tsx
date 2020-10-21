import React, { ChangeEvent, PureComponent } from 'react';
import { LegacyForms } from '@grafana/ui';
import { DataSourcePluginOptionsEditorProps } from '@grafana/data';
import { MyDataSourceOptions, MySecureJsonData } from './types';

const { SecretFormField, FormField } = LegacyForms;

interface Props extends DataSourcePluginOptionsEditorProps<MyDataSourceOptions> {}

interface State {}

export class ConfigEditor extends PureComponent<Props, State> {
  onAuthServerUrlChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onOptionsChange, options } = this.props;
    const jsonData = {
      ...options.jsonData,
      authServerUrl: event.target.value,
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
      clientId: event.target.value,
    };
    onOptionsChange({ ...options, jsonData });
  };

  onContextBrokerUrlChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onOptionsChange, options } = this.props;
    const jsonData = {
      ...options.jsonData,
      contextBrokerUrl: event.target.value,
    };
    onOptionsChange({ ...options, jsonData });
  };

  // Secure field (only sent to the backend)
  onClientSecretChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onOptionsChange, options } = this.props;
    onOptionsChange({
      ...options,
      secureJsonData: {
        clientSecret: event.target.value,
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
            labelWidth={9}
            inputWidth={22}
            onChange={this.onAuthServerUrlChange}
            value={jsonData.authServerUrl || ''}
            placeholder="https://my.sso.org"
          />
        </div>

        <div className="gf-form">
          <FormField
            label="Token endpoint path"
            labelWidth={9}
            inputWidth={22}
            onChange={this.onResourceChange}
            value={jsonData.resource || ''}
            placeholder="/path/to/token/endpoint"
          />
        </div>

        <div className="gf-form">
          <FormField
            label="Client id"
            labelWidth={9}
            inputWidth={22}
            onChange={this.onClientIdChange}
            value={jsonData.clientId || ''}
            tooltip="OAuth2 client id to be used by the plugin"
          />
        </div>

        <div className="gf-form-inline">
          <div className="gf-form">
            <SecretFormField
              isConfigured={(secureJsonFields && secureJsonFields.client_secret) as boolean}
              value={secureJsonData.clientSecret || ''}
              label="Client secret"
              tooltip="OAuth client secret to be used by the plugin"
              placeholder=""
              labelWidth={9}
              inputWidth={22}
              onReset={this.onResetClientSecretKey}
              onChange={this.onClientSecretChange}
            />
          </div>
        </div>

        <div className="gf-form">
          <FormField
            label="NGSI-LD API URL"
            labelWidth={9}
            inputWidth={22}
            onChange={this.onContextBrokerUrlChange}
            value={jsonData.contextBrokerUrl || ''}
            placeholder="https://my.context-brocker.org"
          />
        </div>
      </div>
    );
  }
}
