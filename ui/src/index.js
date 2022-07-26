
import { createRoot } 
  from 'react-dom/client';

import reportWebVitals
  from './reportWebVitals';

import 'bootstrap/dist/css/bootstrap.css';
import './scss/main.scss';

import Main 
  from './app/Main';

import Api
  from './api'

// Alice theme and extension API 
window.Alice = Api;

const root = createRoot(document.getElementById('app'));

root.render(<Main />);

reportWebVitals();
