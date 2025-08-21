import { Suspense, lazy } from 'react';
import { BrowserRouter, Routes, Route } from 'react-router-dom';

const HomePage = lazy(() => import('./pages/Home'));

function App() {
  return (
    <BrowserRouter>
      <Suspense
        fallback={
          <div className="flex items-center justify-center min-h-screen">
            Loading...
          </div>
        }
      >
        <Routes>
          <Route path="/" element={<HomePage />} />
          <Route
            path="*"
            element={
              <div className="flex items-center justify-center min-h-screen">
                <div className="text-center">
                  <h1 className="text-4xl font-bold mb-4">404</h1>
                  <p className="text-lg text-muted-foreground">
                    Page not found
                  </p>
                </div>
              </div>
            }
          />
        </Routes>
      </Suspense>
    </BrowserRouter>
  );
}

export default App;
