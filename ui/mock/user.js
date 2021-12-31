const mockjs= require('mockjs');
const { VUE_APP_APIHOST_MOCK } = process.env;
const ajaxHeadersTokenKey = 'authorization';
const mock = {};

mock[`GET ${VUE_APP_APIHOST_MOCK}/users/profile`] = (req, res) => {
    const headers = req.headers;
    if (headers[ajaxHeadersTokenKey] === 'Bearer admin') {
        res.send({
          code: 0,
          data: {
            id: 1,
            name: 'Admins',
            avatar: '',
            roles: ['admin'],
          },
        });
    } else if (headers[ajaxHeadersTokenKey] === 'Bearer user') {
        res.send({
          code: 0,
          data: {
            id: 2,
            name: 'Users',
            avatar: '',
            roles: ['user'],
          },
        });
    } else if (headers[ajaxHeadersTokenKey] === 'test') {
        res.send({
          code: 0,
          data: {
            id: 3,
            name: 'Tests',
            avatar: '',
            roles: ['test'],
          },
        });
    } else {
        res.send({
          code: 10002,
          data: {headers: req.headers},
          msg: '未登录',
        });
    }

};

mock[`GET ${VUE_APP_APIHOST_MOCK || ''}/users/message`] = (req, res) => {
    res.send({
      code: 0,
      data: mockjs.mock('@integer(0,99)'),
    });
};
  
mock[`POST ${VUE_APP_APIHOST_MOCK || ''}/account/login`] = (req, res) => {
    const { password, username } = req.body;
    const send = { code: 0, data: {}, msg: '' };
    if (username === 'admin' && password === 'password') {
        send['data'] = {
        token: 'admin',
        };
    } else if (username === 'user' && password === 'password') {
        send['data'] = {
        token: 'user',
        };
    } else if (username === 'test' && password === 'password') {
        send['data'] = {
        token: 'test',
        };
    } else {
        send['code'] = 201;
        send['msg'] = 'Wrong username or password';
    }

    res.send(send);
};
  
mock[`POST ${VUE_APP_APIHOST_MOCK || ''}/users/register`] = (req, res) => {
    res.send({
      code: 0,
      data: '',
      msg: '',
    });
};
  

module.exports = {
  ...mock
};