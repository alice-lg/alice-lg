
import { StrictMode }
  from 'react';

import { createRoot } 
  from 'react-dom/client';

import reportWebVitals
  from './reportWebVitals';

import 'bootstrap/dist/css/bootstrap.css';
import './scss/main.scss';

import Alice 
  from './app/Alice';

const root = createRoot(document.getElementById('root'));

root.render(
  <StrictMode>
    <Alice />
  </StrictMode>,
);

reportWebVitals();
