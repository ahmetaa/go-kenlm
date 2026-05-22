#include "wrap.h"

#include "lm/model.hh"
#include "lm/virtual_interface.hh"
#include "util/string_piece.hh"

#include <cstdlib>
#include <cstring>
#include <exception>
#include <sstream>
#include <string>
#include <vector>

namespace {

char* copy_error(const char* msg) {
  size_t n = std::strlen(msg);
  char* out = static_cast<char*>(std::malloc(n + 1));
  if (out == nullptr) return nullptr;
  std::memcpy(out, msg, n + 1);
  return out;
}

}  // namespace

struct kenlm_model {
  lm::base::Model* model;
};

extern "C" {

kenlm_model* kenlm_load(const char* path, char** err_out) {
  try {
    lm::ngram::Config config;
    lm::base::Model* m = lm::ngram::LoadVirtual(path, config);
    kenlm_model* wrap = new kenlm_model;
    wrap->model = m;
    return wrap;
  } catch (const std::exception& e) {
    if (err_out) *err_out = copy_error(e.what());
    return nullptr;
  } catch (...) {
    if (err_out) *err_out = copy_error("unknown error loading kenlm model");
    return nullptr;
  }
}

void kenlm_free(kenlm_model* m) {
  if (m == nullptr) return;
  delete m->model;
  delete m;
}

unsigned int kenlm_order(const kenlm_model* m) {
  return m->model->Order();
}

double kenlm_score_sentence(const kenlm_model* m, const char* tokens) {
  const lm::base::Model* model = m->model;
  const lm::base::Vocabulary& vocab = model->BaseVocabulary();

  const size_t state_size = model->StateSize();
  std::vector<char> a(state_size), b(state_size);
  void* in_state = a.data();
  void* out_state = b.data();

  // Start from begin-of-sentence state (already accounts for <s>).
  model->BeginSentenceWrite(in_state);

  double total = 0.0;
  std::istringstream ss(tokens ? tokens : "");
  std::string token;
  while (ss >> token) {
    lm::WordIndex w = vocab.Index(token);
    total += model->BaseScore(in_state, w, out_state);
    std::swap(in_state, out_state);
  }
  // Terminal </s>.
  total += model->BaseScore(in_state, vocab.EndSentence(), out_state);
  return total;
}

void kenlm_string_free(char* s) {
  std::free(s);
}

}  // extern "C"
