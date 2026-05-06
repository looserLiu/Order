import 'package:flutter_test/flutter_test.dart';
import 'package:order/main.dart';

void main() {
  testWidgets('App starts successfully', (WidgetTester tester) async {
    await tester.pumpWidget(const OrderApp());
    expect(find.text('记账'), findsWidgets);
  });
}