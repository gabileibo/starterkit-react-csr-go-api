import { UserList } from '@/components/UserList';

export default function Home() {
  return (
    <div className="min-h-screen bg-gray-50">
      <div className="container mx-auto px-4 py-8">
        <div className="max-w-6xl mx-auto">
          <div className="mb-8">
            <h1 className="text-3xl font-bold text-gray-900 mb-2">Users</h1>
            <p className="text-gray-600">
              Manage and view all users in the system
            </p>
          </div>

          <div className="bg-white rounded-lg shadow">
            <UserList />
          </div>
        </div>
      </div>
    </div>
  );
}
