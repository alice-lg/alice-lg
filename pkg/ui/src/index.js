
import { StrictMode }
  from 'react';

import { createRoot } 
  from 'react-dom/client';

import reportWebVitals
  from './reportWebVitals';

import 'bootstrap/dist/css/bootstrap.css';
import './scss/main.scss';

import Main 
  from './app/Main';

const root = createRoot(document.getElementById('app'));

root.render(
  <StrictMode><Main /></StrictMode>,
);

reportWebVitals();
