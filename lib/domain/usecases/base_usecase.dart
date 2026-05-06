/// Base class for all Use Cases.
/// Use Cases represent single actions in the domain layer.
abstract class UseCase<Type, Params> {
  /// Execute the use case with given parameters.
  Future<Type> call(Params params);
}

/// Use case with no parameters.
abstract class NoParamsUseCase<Type> {
  /// Execute the use case.
  Future<Type> call();
}
