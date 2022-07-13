
import { createRoot } 
  from 'react-dom/client';

import reportWebVitals
  from './reportWebVitals';

import 'bootstrap/dist/css/bootstrap.css';
import './scss/main.scss';

import Main 
  from './app/Main';

const root = createRoot(document.getElementById('app'));

root.render(<Main />);

reportWebVitals();
