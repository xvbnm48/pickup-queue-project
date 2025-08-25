'use client';

import { useState, useEffect } from 'react';

interface Package {
    id: string;
    order_reference: string;
    driver_code?: string;
    status: 'WAITING' | 'PICKED' | 'HANDED_OVER' | 'EXPIRED';
    created_at: string;
    updated_at: string;
}

interface Stats {
    total: number;
    waiting: number;
    picked: number;
    handed_over: number;
    expired: number;
}

type StatusFilter = 'ALL' | 'WAITING' | 'PICKED' | 'HANDED_OVER' | 'EXPIRED';

export default function PickupDashboard() {
    const [packages, setPackages] = useState<Package[]>([]);
    const [filteredPackages, setFilteredPackages] = useState<Package[]>([]);
    const [stats, setStats] = useState<Stats>({
        total: 0,
        waiting: 0,
        picked: 0,
        handed_over: 0,
        expired: 0,
    });
    const [statusFilter, setStatusFilter] = useState<StatusFilter>('ALL');
    const [showAddModal, setShowAddModal] = useState(false);
    const [showUpdateModal, setShowUpdateModal] = useState(false);
    const [selectedPackage, setSelectedPackage] = useState<Package | null>(null);
    const [selectedNewStatus, setSelectedNewStatus] = useState<string>('');
    const [newPackage, setNewPackage] = useState({
        order_reference: '',
        driver_code: '',
    });

    // Fetch packages and stats
    const fetchData = async () => {
        try {
            // Fetch packages
            const packagesResponse = await fetch('http://localhost:8080/api/v1/packages');
            console.log('Packages Response:', packagesResponse);
            if (packagesResponse.ok) {
                const packagesData = await packagesResponse.json();
                setPackages(packagesData.data || []);
            }

            // Fetch stats
            const statsResponse = await fetch('http://localhost:8080/api/v1/packages/stats');
            console.log('Stats Response:', statsResponse);
            if (statsResponse.ok) {
                const statsData = await statsResponse.json();
                setStats(statsData.data || stats);
            }
        } catch (error) {
            console.error('Error fetching data:', error);
        }
    };

    useEffect(() => {
        fetchData();
        const interval = setInterval(fetchData, 5000); // Refresh every 5 seconds
        return () => clearInterval(interval);
    }, []);

    // Filter packages based on status
    useEffect(() => {
        if (statusFilter === 'ALL') {
            setFilteredPackages(packages);
        } else {
            setFilteredPackages(packages.filter(pkg => pkg.status === statusFilter));
        }
    }, [packages, statusFilter]);

    // Create package
    const handleCreatePackage = async (e: React.FormEvent) => {
        e.preventDefault();
        try {
            const response = await fetch('http://localhost:8080/api/v1/packages', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(newPackage),
            });

            if (response.ok) {
                setShowAddModal(false);
                setNewPackage({ order_reference: '', driver_code: '' });
                fetchData();
            }
        } catch (error) {
            console.error('Error creating package:', error);
        }
    };

    // Update package status
    const handleSaveStatusChange = async () => {
        if (!selectedPackage || !selectedNewStatus) return;

        try {
            const response = await fetch(`http://localhost:8080/api/v1/packages/${selectedPackage.id}/status`, {
                method: 'PATCH',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ status: selectedNewStatus }),
            });

            if (response.ok) {
                setShowUpdateModal(false);
                setSelectedPackage(null);
                setSelectedNewStatus('');
                fetchData();
            }
        } catch (error) {
            console.error('Error updating package status:', error);
        }
    };

    const getStatusColor = (status: string) => {
        switch (status) {
            case 'WAITING':
                return 'bg-yellow-100 text-yellow-800';
            case 'PICKED':
                return 'bg-blue-100 text-blue-800';
            case 'HANDED_OVER':
                return 'bg-green-100 text-green-800';
            case 'EXPIRED':
                return 'bg-red-100 text-red-800';
            default:
                return 'bg-gray-100 text-gray-800';
        }
    };

    const getValidNextStatuses = (currentStatus: string) => {
        switch (currentStatus) {
            case 'WAITING':
                return ['PICKED'];
            case 'PICKED':
                return ['HANDED_OVER'];
            default:
                return [];
        }
    };

    return (
        <div className="min-h-screen bg-gray-50 p-6">
            <div className="max-w-7xl mx-auto">
                {/* Header */}
                <div className="mb-8">
                    <h1 className="text-3xl font-bold text-gray-900 mb-2">Pickup Queue</h1>

                    {/* Status Filters */}
                    <div className="flex gap-2 mb-4">
                        {(['ALL', 'WAITING', 'PICKED', 'HANDED_OVER', 'EXPIRED'] as StatusFilter[]).map((status) => (
                            <button
                                key={status}
                                onClick={() => setStatusFilter(status)}
                                className={`px-4 py-2 rounded-full text-sm font-medium transition-colors ${statusFilter === status
                                    ? 'bg-blue-600 text-white'
                                    : 'bg-white text-gray-700 hover:bg-gray-50 border border-gray-300'
                                    }`}
                            >
                                {status}
                            </button>
                        ))}
                    </div>

                    {/* Add Package Button */}
                    <button
                        onClick={() => setShowAddModal(true)}
                        className="bg-green-600 text-white px-4 py-2 rounded-full hover:bg-green-700 transition-colors"
                    >
                        + Add Package
                    </button>
                </div>

                {/* Stats Cards */}
                <div className="grid grid-cols-2 md:grid-cols-5 gap-4 mb-8">
                    <div className="bg-white p-4 rounded-lg border border-gray-200">
                        <div className="text-2xl font-bold text-gray-900">{stats.total}</div>
                        <div className="text-sm text-gray-600">Total</div>
                    </div>
                    <div className="bg-white p-4 rounded-lg border border-gray-200">
                        <div className="text-2xl font-bold text-yellow-600">{stats.waiting}</div>
                        <div className="text-sm text-gray-600">WAITING</div>
                    </div>
                    <div className="bg-white p-4 rounded-lg border border-gray-200">
                        <div className="text-2xl font-bold text-blue-600">{stats.picked}</div>
                        <div className="text-sm text-gray-600">PICKED</div>
                    </div>
                    <div className="bg-white p-4 rounded-lg border border-gray-200">
                        <div className="text-2xl font-bold text-green-600">{stats.handed_over}</div>
                        <div className="text-sm text-gray-600">HANDED_OVER</div>
                    </div>
                    <div className="bg-white p-4 rounded-lg border border-gray-200">
                        <div className="text-2xl font-bold text-red-600">{stats.expired}</div>
                        <div className="text-sm text-gray-600">EXPIRED</div>
                    </div>
                </div>

                {/* Packages Table */}
                <div className="bg-white rounded-lg border border-gray-200 overflow-hidden">
                    <div className="overflow-x-auto">
                        <table className="w-full">
                            <thead className="bg-gray-50">
                                <tr>
                                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                        Package ID
                                    </th>
                                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                        Order Ref
                                    </th>
                                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                        Driver
                                    </th>
                                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                        Status
                                    </th>
                                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                        Actions
                                    </th>
                                </tr>
                            </thead>
                            <tbody className="bg-white divide-y divide-gray-200">
                                {filteredPackages.map((pkg) => (
                                    <tr key={pkg.id} className="hover:bg-gray-50">
                                        <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                                            {pkg.id.substring(0, 8)}...
                                        </td>
                                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                                            {pkg.order_reference}
                                        </td>
                                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                                            {pkg.driver_code || '-'}
                                        </td>
                                        <td className="px-6 py-4 whitespace-nowrap">
                                            <span className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${getStatusColor(pkg.status)}`}>
                                                {pkg.status}
                                            </span>
                                        </td>
                                        <td className="px-6 py-4 whitespace-nowrap text-sm">
                                            {getValidNextStatuses(pkg.status).length > 0 && (
                                                <button
                                                    onClick={() => {
                                                        setSelectedPackage(pkg);
                                                        setSelectedNewStatus('');
                                                        setShowUpdateModal(true);
                                                    }}
                                                    className="bg-blue-600 text-white px-3 py-1 rounded-full hover:bg-blue-700 transition-colors text-sm font-medium"
                                                >
                                                    Update
                                                </button>
                                            )}
                                        </td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                    </div>
                </div>

                {/* Add Package Modal */}
                {showAddModal && (
                    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
                        <div className="bg-white p-6 rounded-lg w-96 shadow-xl">
                            <h2 className="text-xl font-semibold mb-6 text-gray-900">Add New Package</h2>
                            <form onSubmit={handleCreatePackage}>
                                <div className="mb-4">
                                    <label className="block text-sm font-medium text-gray-700 mb-2">
                                        Order Reference *
                                    </label>
                                    <input
                                        type="text"
                                        value={newPackage.order_reference}
                                        onChange={(e) => setNewPackage({ ...newPackage, order_reference: e.target.value })}
                                        className="w-full p-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 text-gray-900 placeholder-gray-400"
                                        placeholder="Enter order reference"
                                        required
                                    />
                                </div>
                                <div className="mb-6">
                                    <label className="block text-sm font-medium text-gray-700 mb-2">
                                        Driver Code (opt)
                                    </label>
                                    <input
                                        type="text"
                                        value={newPackage.driver_code}
                                        onChange={(e) => setNewPackage({ ...newPackage, driver_code: e.target.value })}
                                        className="w-full p-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 text-gray-900 placeholder-gray-400"
                                        placeholder="Enter driver code (optional)"
                                    />
                                </div>
                                <div className="flex gap-3">
                                    <button
                                        type="submit"
                                        className="flex-1 bg-green-600 text-white py-3 px-4 rounded-full hover:bg-green-700 transition-colors font-medium"
                                    >
                                        Create
                                    </button>
                                    <button
                                        type="button"
                                        onClick={() => {
                                            setShowAddModal(false);
                                            setNewPackage({ order_reference: '', driver_code: '' });
                                        }}
                                        className="flex-1 bg-gray-300 text-gray-700 py-3 px-4 rounded-full hover:bg-gray-400 transition-colors font-medium"
                                    >
                                        Cancel
                                    </button>
                                </div>
                            </form>
                        </div>
                    </div>
                )}

                {/* Update Status Modal */}
                {showUpdateModal && selectedPackage && (
                    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
                        <div className="bg-white p-6 rounded-lg w-96 shadow-xl">
                            <h2 className="text-xl font-semibold mb-6 text-gray-900">Update Package #{selectedPackage.id.substring(0, 8)}...</h2>
                            <div className="mb-6">
                                <p className="text-sm text-gray-700 mb-4">
                                    Current status: <span className="font-semibold text-gray-900">{selectedPackage.status}</span>
                                </p>
                                <label className="block text-sm font-medium text-gray-700 mb-2">
                                    New status:
                                </label>
                                <select
                                    value={selectedNewStatus}
                                    onChange={(e) => setSelectedNewStatus(e.target.value)}
                                    className="w-full p-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 text-gray-900"
                                >
                                    <option value="">Select new status</option>
                                    {getValidNextStatuses(selectedPackage.status).map((status) => (
                                        <option key={status} value={status}>
                                            {status}
                                        </option>
                                    ))}
                                </select>
                            </div>
                            <div className="flex gap-3 mt-6">
                                <button
                                    onClick={handleSaveStatusChange}
                                    disabled={!selectedNewStatus}
                                    className="flex-1 bg-green-600 text-white py-3 px-4 rounded-full hover:bg-green-700 disabled:bg-gray-300 disabled:cursor-not-allowed transition-colors font-medium"
                                >
                                    Save Change
                                </button>
                                <button
                                    onClick={() => {
                                        setShowUpdateModal(false);
                                        setSelectedPackage(null);
                                        setSelectedNewStatus('');
                                    }}
                                    className="flex-1 bg-gray-300 text-gray-700 py-3 px-4 rounded-full hover:bg-gray-400 transition-colors font-medium"
                                >
                                    Cancel
                                </button>
                            </div>
                        </div>
                    </div>
                )}
            </div>
        </div>
    );
}
