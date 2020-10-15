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

  onApiResourceChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onOptionsChange, options } = this.props;
    const jsonData = {
      ...options.jsonData,
      apiResource: event.target.value,
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
            label="token URL"
            labelWidth={7}
            inputWidth={22}
            onChange={this.onTokenUrlChange}
            value={jsonData.tokenUrl || ''}
            placeholder="url to get token"
          />
        </div>

        <div className="gf-form">
          <FormField
            label="resource URL"
            labelWidth={7}
            inputWidth={22}
            onChange={this.onResourceChange}
            value={jsonData.resource || ''}
            placeholder="ressources to complete URL"
          />
        </div>

        <div className="gf-form">
          <FormField
            label="client_id"
            labelWidth={7}
            inputWidth={22}
            onChange={this.onClientIdChange}
            value={jsonData.client_id || ''}
            placeholder="client_id to get token"
          />
        </div>

        <div className="gf-form-inline">
          <div className="gf-form">
            <SecretFormField
              isConfigured={(secureJsonFields && secureJsonFields.client_secret) as boolean}
              value={secureJsonData.client_secret || ''}
              label="client_secret"
              placeholder="client_secret to get token"
              labelWidth={7}
              inputWidth={22}
              onReset={this.onResetClientSecretKey}
              onChange={this.onClientSecretChange}
            />
          </div>
        </div>

        <div className="gf-form">
          <FormField
            label="api URL"
            labelWidth={7}
            inputWidth={22}
            onChange={this.onApiUrlChange}
            value={jsonData.apiUrl || ''}
            placeholder="api URL to make http request"
          />
        </div>

        <div className="gf-form">
          <FormField
            label="api resource URL"
            labelWidth={7}
            inputWidth={22}
            onChange={this.onApiResourceChange}
            value={jsonData.apiResource || ''}
            placeholder="resource to complete api URL"
          />
        </div>
      </div>
    );
  }
}
