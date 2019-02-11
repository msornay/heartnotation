import React from 'react';
import Enzyme, { mount, ReactWrapper } from 'enzyme';
import Adapter from 'enzyme-adapter-react-16';
import Dashboard, { Annotation } from '../pages/Dashboard';
import { Table, Input } from 'antd';
import mock from './mocks/annotations.json';

Enzyme.configure({ adapter: new Adapter() });

const mockData: Annotation[] = JSON.parse(JSON.stringify(mock));

const mockPromise = () => Promise.resolve(mockData);

describe('Dashboard page', () => {
  let mountedPage: ReactWrapper;

  beforeEach(() => {
    mountedPage = mount(<Dashboard getAnnotations={mockPromise} />);
  });

  it('renders', () => {
    expect(mountedPage).not.toBeFalsy();
  });

  describe('after API call', async () => {
    mockPromise().then(() => {
      it('state have correct length', () => {
        expect(mountedPage.state('initialAnnotations')).toHaveLength(
          mockData.length
        );
      });

      it('have correct number of rows', () => {
        expect(mountedPage.find(Table).props().dataSource).toHaveLength(
          mockData.length
        );
      });

      it('filters on id column', () => {
        mountedPage
          .findWhere(p => p.hasClass('search_id'))
          .find(Input)
          .simulate('change', { target: { value: '1' } });
        expect(mountedPage.find(Table).props().dataSource).toHaveLength(1);
      });

      it('filters on signal id column', () => {
        mountedPage
          .findWhere(p => p.hasClass('search_signal_id'))
          .find(Input)
          .simulate('change', { target: { value: '1' } });
        expect(mountedPage.find(Table).props().dataSource).toHaveLength(4);
      });

      it('filters on name column', () => {
        mountedPage
          .findWhere(p => p.hasClass('search_name'))
          .find(Input)
          .simulate('change', { target: { value: 'base' } });
        expect(mountedPage.find(Table).props().dataSource).toHaveLength(2);
      });

      it('reset filters after cleaning fields', () => {
        mountedPage
          .findWhere(p => p.hasClass('search_name'))
          .find(Input)
          .simulate('change', { target: { value: 'base' } });
        expect(mountedPage.find(Table).props().dataSource).toHaveLength(2);
        mountedPage
          .findWhere(p => p.hasClass('search_name'))
          .find(Input)
          .simulate('change', { target: { value: '' } });
        expect(mountedPage.find(Table).props().dataSource).toHaveLength(4);
      });
    });
  });
});