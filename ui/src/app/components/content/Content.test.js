import { useEffect } from 'react';

import { render, screen } from '@testing-library/react';

import Content from 'app/components/content/Content';
import { ContentProvider } from 'app/context/content';
import { updateContent } from 'api';

test("render Content with test context", () => {

  const App = () => {
    useEffect(() => {
      updateContent({"cid": "test123"});
    });

    return (
      <ContentProvider>
        <Content id="cid" />
      </ContentProvider>
    );
  };

  render(<App />);
  expect(screen.queryByText("test123")).not.toBe(null);
});

