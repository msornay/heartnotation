import React, { Component } from 'react';
import { Row, Col, Icon, Switch, Button, Steps, Alert } from 'antd';
import { Annotation, StatusInserter, Status, api } from '../../utils';
import { RouteComponentProps, withRouter } from 'react-router';
import { withAuth, AuthProps } from '../../utils/auth';
import ChatDrawerAnnotation from '../chatAnnotation/ChatDrawerAnnotation';

interface State {
  stepProcess: number;
  mode: 'Navigation' | 'Annotation';
  error?: string;
}

interface Props extends RouteComponentProps, AuthProps {
  annotation: Annotation;
  onToggle: (state: boolean) => void;
}

interface PropsButton extends AuthProps {
  conditionnal_id: number;
  annotation: Annotation;
  handleSubmit: (s: StatusInserter) => void;
}

const ValidateButton = (props: PropsButton) => {
  const { user } = props;
  return props.user.role.name === 'Gestionnaire' ? (
    <Button
      type='primary'
      icon='check-circle'
      size='large'
      onClick={() => {
        props.handleSubmit({
          status: 5,
          id: props.annotation.id
        });
      }}
    >
      Validate
    </Button>
  ) : null;
};

const InvalidateButton = (props: PropsButton) => {
  const { user } = props;
  return props.user.role.name === 'Gestionnaire' ? (
    <Button
      type='danger'
      icon='close-circle'
      size='large'
      onClick={() => {
        props.handleSubmit({
          status: 3,
          id: props.annotation.id
        });
      }}
    >
      Invalidate
    </Button>
  ) : null;
};

const CompleteButton = (props: PropsButton) => {
  const { user } = props;
  return props.user.role.name === 'Annotateur' ? (
    <Button
      type='default'
      icon='check-circle'
      size='large'
      onClick={() => {
        props.handleSubmit({
          status: 4,
          id: props.annotation.id
        });
      }}
    >
      Complete
    </Button>
  ) : null;
};

const ConditionalButton = (props: PropsButton) => {
  const { conditionnal_id } = props;
  if (conditionnal_id === 0) {
    return <CompleteButton {...props} />;
  } else if (conditionnal_id === 1) {
    return (
      <>
        <Col span={12}>
          <InvalidateButton key={1} {...props} />
        </Col>
        <Col span={12}>
          <ValidateButton key={2} {...props} />
        </Col>
      </>
    );
  } else if (conditionnal_id === 2) {
    return null;
  }
  return null;
};

class HeaderSignalAnnotation extends Component<Props, State> {
  constructor(props: Props) {
    super(props);
    let step = -1;
    if (props.annotation.status) {
      props.annotation.status.sort(
        (s1: Status, s2: Status) => s2.date.getTime() - s1.date.getTime()
      );
      switch (props.annotation.status[0].enum_status.name) {
        case 'ASSIGNED':
        case 'IN_PROCESS':
          step = 0;
          break;
        case 'COMPLETED':
          step = 1;
          break;
        case 'VALIDATED':
          step = 2;
          break;
        default:
          break;
      }
      this.state = {
        stepProcess: step,
        mode: 'Navigation'
      };
    }
  }

  public handleSubmit = (s: StatusInserter) => {
    api.sendStatus(s).then(() => {
      this.props.history.push('/');
    });
    // .catch(error => {
    //   this.setState({ error });
    // });
  }

  private handleToggle = (toggle: boolean) => {
    this.setState({ mode: toggle ? 'Annotation' : 'Navigation' }, () =>
      this.props.onToggle(toggle)
    );
  }

  public render() {
    const { Step } = Steps;
    const { annotation, user } = this.props;
    const { stepProcess, error, mode } = this.state;
    console.log(stepProcess);
    return [
      <Row
        key={1}
        type='flex'
        className='signal-header'
        align='middle'
        justify='space-between'
      >
        <Col span={4}>
          <Switch
            checkedChildren={<Icon type='check' />}
            unCheckedChildren={<Icon type='close' />}
            defaultChecked={true}
          />
          Display Leads
        </Col>
        {user.role.name === 'Annotateur' && stepProcess === 0 && (
          <Col span={4}>
            {mode} Mode <Switch onChange={this.handleToggle} />
          </Col>
        )}
        <Col span={8}>
          <Steps
            style={{ paddingTop: 30 }}
            progressDot={true}
            current={stepProcess}
            size='default'
          >
            <Step title='In Progress' />
            <Step title='Completed' />
            <Step title='Validated' />
          </Steps>
        </Col>
        <Col offset={1} span={3}>
          <ChatDrawerAnnotation />
        </Col>
        <Col span={4}>
          <Row type='flex' align='middle' justify='end'>
            <ConditionalButton
              conditionnal_id={stepProcess}
              annotation={annotation}
              handleSubmit={this.handleSubmit}
              user={user}
            />
            {error && <Alert message={error} type='error' />}
          </Row>
        </Col>
      </Row>
    ];
  }
}

export default withRouter(withAuth(HeaderSignalAnnotation));
