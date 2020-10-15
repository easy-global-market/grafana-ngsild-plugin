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
            label="Token URL"
            labelWidth={7}
            inputWidth={22}
            onChange={this.onTokenUrlChange}
            value={jsonData.tokenUrl || ''}
            placeholder="http://localhost:3306"
            tooltip="URL to get token (without path)"
          />
        </div>

        <div className="gf-form">
          <FormField
            label="Path"
            labelWidth={7}
            inputWidth={22}
            onChange={this.onResourceChange}
            value={jsonData.resource || ''}
            placeholder="/auth/.../token"
            tooltip="Path to complete token URL"
          />
        </div>

        <div className="gf-form">
          <FormField
            label="Client_id"
            labelWidth={7}
            inputWidth={22}
            onChange={this.onClientIdChange}
            value={jsonData.client_id || ''}
            placeholder="id"
            tooltip="client_id to get token"
          />
        </div>

        <div className="gf-form-inline">
          <div className="gf-form">
            <SecretFormField
              isConfigured={(secureJsonFields && secureJsonFields.client_secret) as boolean}
              value={secureJsonData.client_secret || ''}
              label="Client_secret"
              placeholder="secret"
              tooltip="client_secret to get token"
              labelWidth={7}
              inputWidth={22}
              onReset={this.onResetClientSecretKey}
              onChange={this.onClientSecretChange}
            />
          </div>
        </div>

        <div className="gf-form">
          <FormField
            label="Api URL"
            labelWidth={7}
            inputWidth={22}
            onChange={this.onApiUrlChange}
            value={jsonData.apiUrl || ''}
            placeholder="http://localhost:3307"
            tooltip="Api URL to make http request (without path)"
          />
        </div>

        <div className="gf-form">
          <FormField
            label="Api path"
            labelWidth={7}
            inputWidth={22}
            onChange={this.onApiResourceChange}
            value={jsonData.apiResource || ''}
            placeholder="/v1/.../entities/"
            tooltip="Api path to complete api URL"
          />
        </div>
      </div>
    );
  }
}
