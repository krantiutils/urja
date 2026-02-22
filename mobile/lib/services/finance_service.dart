import '../models/due.dart';
import '../models/transaction.dart';
import 'api_client.dart';

class FinanceService {
  final ApiClient _api;

  FinanceService(this._api);

  // Dues
  Future<({List<Due> data, int total})> getDues(
    String orgId, {
    String? status,
    String? search,
    int limit = 20,
    int offset = 0,
  }) async {
    final response = await _api.get('/orgs/$orgId/dues', queryParameters: {
      'limit': limit,
      'offset': offset,
      if (status != null && status.isNotEmpty) 'status': status,
      if (search != null && search.isNotEmpty) 'search': search,
    });
    final list = (response.data['data'] as List<dynamic>? ?? [])
        .map((e) => Due.fromJson(e as Map<String, dynamic>))
        .toList();
    return (data: list, total: response.data['total'] as int? ?? list.length);
  }

  Future<void> payDue(String orgId, String dueId,
      {required String amount,
      required String paymentMethod,
      String? paymentReference}) async {
    await _api.post('/orgs/$orgId/dues/$dueId/pay', data: {
      'amount': amount,
      'payment_method': paymentMethod,
      if (paymentReference != null) 'payment_reference': paymentReference,
    });
  }

  // Accounts / Transactions
  Future<({List<Transaction> data, int total})> getTransactions(
    String orgId, {
    String? category,
    String? transactionType,
    String? from,
    String? to,
    int limit = 20,
    int offset = 0,
  }) async {
    final response = await _api.get('/orgs/$orgId/accounts', queryParameters: {
      'limit': limit,
      'offset': offset,
      if (category != null && category.isNotEmpty) 'category': category,
      if (transactionType != null && transactionType.isNotEmpty)
        'transaction_type': transactionType,
      if (from != null) 'from': from,
      if (to != null) 'to': to,
    });
    final list = (response.data['data'] as List<dynamic>? ?? [])
        .map((e) => Transaction.fromJson(e as Map<String, dynamic>))
        .toList();
    return (data: list, total: response.data['total'] as int? ?? list.length);
  }

  Future<void> addTransaction(
      String orgId, Map<String, dynamic> data) async {
    await _api.post('/orgs/$orgId/accounts', data: data);
  }

  Future<void> updateTransaction(
      String orgId, String id, Map<String, dynamic> data) async {
    await _api.put('/orgs/$orgId/accounts/$id', data: data);
  }

  Future<void> deleteTransaction(String orgId, String id) async {
    await _api.deleteRequest('/orgs/$orgId/accounts/$id');
  }

  Future<AccountsSummary> getSummary(String orgId,
      {String? from, String? to}) async {
    final response =
        await _api.get('/orgs/$orgId/accounts/summary', queryParameters: {
      if (from != null) 'from': from,
      if (to != null) 'to': to,
    });
    return AccountsSummary.fromJson(response.data);
  }
}
